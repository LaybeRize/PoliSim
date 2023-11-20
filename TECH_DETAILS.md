# PoliSim - A Political Simulation Website

The project has a [public Docker repo](https://hub.docker.com/r/layberize/polisim) where the newest version can be found as a container. An example docker-compose file
can be found in the files of the projects. All needed environment parameters to start up the container are listed here:

```env
DB_NAME=prod
DB_PASSWORD=root
DB_USER=postgres
DB_ADRESS=DB adress with port
ADDRESS=host adress with port
COOKIE_KEY=your super secret key here
INIT_NAME=Root Account Name
INIT_USERNAME=Root Account Username
INIT_PASSWORD=Root Account password
CORS_URL=Your URL needed for CORS
LANG=your language (most likely DE or EN)
```

## Architecture

**PoliSim** is build on 5 layers:
1. Database Layer
2. Data Abstraction Layer
3. Data Validation/Processing Layer
4. HTML Composition Layer
5. HTML Serving Layer

**PoliSim** realises the first Layer with GORM. The second layer uses the GORM modells and queries to extract only the needed
data for the request, thereby minimizing the data flow between database and program. This layer is completely encapsulated in the 
extraction package. The third layer takes in the data from the HTML-Request, process it, requests the needed complimentary 
data from the Data Abstraction Layer and gives back the state of validity of the request to the HTML Serving Layer which 
then makes a request to the HTML Composition Layer to generate the appropiate HTML to serve back to the client.

The HTML requests and responses are received and served with go-chi. The HTML is extended with Hyperscript and HTMX to
minimize the needed HTML send back to the client.

## HTML Serving

HTML is build with the component builder which is inspired by [gomponents](https://github.com/maragudk/gomponents).
It uses the idea, but simplifies it down to only four functions, which makes shortens the code
quite significantly. It also uses a different approach when validating if the function is an attribute or not.
The HTML Composition Layer uses this builder to compose the HTML snipptes for the 
HTML Serving Layer.

## Language and Configuration

For most parameter, which should be easily adaptable to a new simulation, there exists a field in the config.json in the "resources" folder.
If you want to support a different language, you have to add a file with your language named
{your two letter shorthand for your language}.json in the resources folder and change your LANG env variable to that shorthand.
Theoretically any name can be chosen for the new file, but because the LANG env variable is used for the html lang attribute too, it
is advisable to use the official shorthand.
This repo accepts pull requests for any new language support. We as the maintainer will inform you of any new translations needed.

## Language error checking

As there is only currently support for german, anyone how speaks german is welcome to check the DE.json file for spelling or grammatical 
errors.  
Support for english is coming, and Spanish/French would be nice, but for these we would need someone that would be willing to open a 
pull request.

# Thanks

Special Thanks to hypermedia to bring me back to writing a good server side application and [a certain JSON extension](https://github.com/Emtyloc/json-enc-custom) for
letting me create an easier voting backend. Also thanks to gomponents (mentioned above) for the inspiration to writing html in go itself.