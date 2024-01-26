import os.path
import re
import xml.etree.ElementTree
from enum import Enum


class Mode(Enum):
    WRITE_START = 1
    SKIP_LINES = 2
    REPLACE_SCHEMA_CREATION = 3
    WRITE_END = 4


go_types_to_db_types = {"int64": "BIGINT",
                        "int8": "SMALLINT",
                        "int": "INT",
                        "bool": "BOOLEAN",
                        "string": "TEXT",
                        "time.Time": "TIMESTAMP WITH TIME ZONE",
                        }


class DBField:
    def __init__(self, name: str, alias: str, go_type: str, extra: dict[str, str], reference: str = ""):
        self.name = name
        self.alias = alias
        if go_type in go_types_to_db_types:
            self.db_type: str = go_types_to_db_types[go_type]
        else:
            self.db_type = extra["dbType"]
        self.go_type = go_type
        self.extra_info = extra
        self.reference = reference


class QueryFunction:
    def __init__(self, name: str, input_parameter: str, return_value: str, query: str, single: bool, single_item: str):
        self.name = name
        self.input_parameter = input_parameter
        self.return_value = return_value
        self.query = query
        self.single = single
        self.single_item = single_item


struct_map = {}
tables = {}
join_tables = {}
column_lookup = {}
queries = {}
star_column_list = {}


def generate_column_text(field_definition: DBField) -> str:
    result = "\n   " + field_definition.name
    if "autoIncrement" in field_definition.extra_info:
        if field_definition.db_type == "INT" or field_definition.db_type == "SMALLINT":
            result += " SERIAL"
        if field_definition.db_type == "BIGINT":
            result += " BIGSERIAL"
    else:
        result += " " + field_definition.db_type
    if "primaryKey" in field_definition.extra_info:
        result += " PRIMARY KEY"
    return result + ","


def create_creation_text() -> str:
    creation_string = ""
    for entry in tables:
        creation_string += "CREATE TABLE IF NOT EXISTS " + entry + " ("
        for element in tables[entry]:
            creation_string += generate_column_text(element)
        creation_string = creation_string[:-1] + "\n);\n\n"

    for entry in join_tables:
        creation_string += "CREATE TABLE IF NOT EXISTS " + entry + " ("
        list_key = ""
        forgein_key_list = ""
        for element in join_tables[entry]:
            creation_string += "\n   " + element.name + " " + element.db_type + " NOT NULL,"
            list_key += element.name + ", "
            forgein_key_list += ("\n   FOREIGN KEY (" + element.name + ") REFERENCES " + element.reference +
                                 " (" + element.name + "),")
        creation_string += "\n   PRIMARY KEY (" + list_key[:-2] + ")," + forgein_key_list
        creation_string = creation_string[:-1] + "\n);\n\n"
    return creation_string


def get_join_child(item) -> DBField:
    temp = column_lookup[item.attrib["table"]][item.attrib["alias"]]
    return DBField(temp.name, temp.alias, temp.go_type, temp.extra_info, item.attrib["table"])


def rewrite_database_go_file():
    with open("..\\database\\database.go", "r+") as go_database:
        result = ""

        mode = Mode.WRITE_START
        for line in go_database.readlines():
            if mode == Mode.WRITE_START:
                if line.startswith("var schema = `"):
                    mode = Mode.SKIP_LINES
                result += line
            if mode == Mode.SKIP_LINES:
                if line.startswith("`"):
                    mode = Mode.REPLACE_SCHEMA_CREATION
            if mode == Mode.REPLACE_SCHEMA_CREATION:
                result += create_creation_text()
                mode = Mode.WRITE_END
            if mode == Mode.WRITE_END:
                result += line

        go_database.seek(0)
        go_database.write(result)
        go_database.truncate()


def get_query_from_node(xml_element) -> QueryFunction:
    ret_val = ""
    query_str = ""
    single = True
    single_item = ""

    for xml_child in xml_element:
        if xml_child.tag == "statement":
            query_str = xml_child.text

        elif xml_child.tag == "return":
            ret_val = xml_child.text

            if xml_child.attrib["amount"] == "multiple":
                single = False
            elif xml_child.attrib["amount"] != "single":
                raise Exception("no valid amount for return value in query with the name: " +
                                xml_element.attrib["name"])
            if "singleVersion" in xml_child.attrib:
                single_item = xml_child.attrib["singleVersion"]

    return QueryFunction(xml_element.attrib["name"], xml_element.attrib["parameter"], ret_val,
                         reslove_query_parameter(query_str), single, single_item)


def reslove_query_parameter(input_query: str) -> str:
    result = input_query
    for table_name in tables:
        if "*_" + table_name in input_query:
            result = str.replace(result, "*_" + table_name, star_column_list[table_name])
        if table_name + "." in input_query:
            for ele in tables[table_name]:
                result = str.replace(result, table_name + "." + ele.alias, table_name + "." + ele.name)
    return result


function_single_str = """
    result := %s{}
    row, err := DB.NamedQuery(`%s`, map[string]interface{}{%s})
    if err != nil {
        return result, err
    }
    row.Next()
    err = row.StructScan(&result)
    return result, err
"""

function_multiple_str = """
    result := make(%s, 0)
    rows, err := DB.NamedQuery(`%s`, map[string]interface{}{%s})
    if err != nil {
        return result, err
    }
    pos := 0
    for rows.Next() {
        result = append(result, %s{})
        err = rows.StructScan(&result[pos])
        if err != nil {
            return result, err
        }
        pos++
    }
    return result, err
"""


def create_query_go_files():
    base_path = "..\\database\\"
    for query_key in queries:
        if not os.path.exists(base_path + query_key):
            with open(base_path + query_key, 'w') as _:
                pass
        with open(base_path + query_key, "r+") as go_file:
            result = "package database\n\n"
            for query_ele in queries[query_key]:
                result += generate_query_function_from_element(query_ele) + "\n"

            go_file.seek(0)
            go_file.write(result)
            go_file.truncate()


def get_map_string_from_parameter(input_parameter: str) -> str:
    first_list = input_parameter.split(",")
    second_list = []
    for ele in first_list:
        second_list.append(ele.split()[0])
    result = ""
    for ele in second_list:
        result += "\""+ele+"\": " + ele + ",\n        "
    return result[:-10]


def generate_query_function_from_element(element: QueryFunction) -> str:
    map_string = get_map_string_from_parameter(element.input_parameter)
    regex = r".{,120} "
    subst = "\\g<0>\\n                                     "
    query_string = re.sub(regex, subst, element.query, 0, re.MULTILINE)

    if element.single:
        result = function_single_str % (element.return_value, query_string, map_string)
    else:
        result = function_multiple_str % (element.return_value, query_string, map_string, element.single_item)
    return ("func " + element.name + "(" + element.input_parameter + ") " + "(" + element.return_value +
            ", error) {" + result + "}")


def update_star_columns():
    for table_name in tables:
        table_column_list = ""
        for ele in tables[table_name]:
            table_column_list += table_name + "." + ele.name + ", "
        star_column_list[table_name] = table_column_list[:-2]


def set_struct_info(xml_child):
    if "struct" in xml_child.attrib:
        struct_map[xml_child.attrib["struct"]] = xml_child.attrib["name"]


def generate_structs_form_map() -> str:
    result = ""
    for entry in struct_map:
        result += "type " + entry + " struct {\n"
        if struct_map[entry] in tables:
            array = tables[struct_map[entry]]
        else:
            array = join_tables[struct_map[entry]]
        for ele in array:
            result += f"   {ele.alias} {ele.go_type} `db:\"{ele.name}\"`\n"
        result += "}\n\n"
    return result


def create_data_structs():
    path = "..\\database\\generalStruct.go"
    if not os.path.exists(path):
        with open(path, 'w') as _:
            pass
    with open(path, "r+") as go_database:
        result = "package database\n\nimport \"time\"\n\n"
        result += generate_structs_form_map()

        go_database.seek(0)
        go_database.write(result)
        go_database.truncate()


if __name__ == '__main__':
    tree = xml.etree.ElementTree.parse("format.xml")
    root = tree.getroot()
    for child in root:
        if child.tag == "table":
            column_lookup[child.attrib["name"]] = {}
            columns = []
            for column in child:
                field = DBField(column.tag, column.attrib["alias"], column.text, column.attrib)
                columns.append(field)
                column_lookup[child.attrib["name"]][column.attrib["alias"]] = field
            tables[child.attrib["name"]] = columns
            set_struct_info(child)

        elif child.tag == "jointable":
            column_lookup[child.attrib["name"]] = {}
            columns = []
            for column in child:
                field = get_join_child(column)
                columns.append(field)
                column_lookup[child.attrib["name"]][field.name] = field
            join_tables[child.attrib["name"]] = columns
            set_struct_info(child)

    update_star_columns()

    for child in root:
        if child.tag == "query":
            corresponding_file = "queries.go"
            if "fileName" in child.attrib:
                corresponding_file = child.attrib["fileName"] + "Queries.go"
            if corresponding_file not in queries:
                queries[corresponding_file] = []
            queries[corresponding_file].append(get_query_from_node(child))

    rewrite_database_go_file()
    create_query_go_files()
    create_data_structs()
