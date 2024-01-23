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
    def __init__(self, name: str, input_parameter: str, return_value: str, query: str, single: bool,
                 use_pointer: bool = False):
        self.name = name
        self.input_parameter = input_parameter
        self.return_value = return_value
        self.query = query
        self.single = single
        self.use_pointer = use_pointer


tables = {}
join_tables = {}
column_lookup = {}
queries = {}


def generate_column_text(field_definition: DBField, add_not_null: bool = False) -> str:
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


def get_query_from_node(child):
    pass


def reslove_query_parameter(input_query: str) -> str:
    pass


if __name__ == '__main__':
    tree = xml.etree.ElementTree.parse("format.xml")
    root = tree.getroot()
    for child in root:
        if child.tag == "table":
            column_lookup[child.attrib["name"]] = {}
            columns = []
            for column in child:
                field = DBField(column.tag, column.attrib["name"], column.text, column.attrib)
                columns.append(field)
                column_lookup[child.attrib["name"]][column.attrib["name"]] = field
            tables[child.attrib["name"]] = columns

        if child.tag == "jointable":
            columns = []
            for column in child:
                field = get_join_child(column)
                columns.append(field)
                column_lookup[child.attrib["name"]][field.name] = field
            join_tables[child.attrib["name"]] = columns
        if child.tag == "query":
            corresponding_file = "queries.go"
            if "fileName" in child.attrib:
                corresponding_file = child.attrib["fileName"] + "Queries.go"
            if corresponding_file not in queries:
                queries[corresponding_file] = []
            queries[corresponding_file].append(get_query_from_node(child))



