# PoliSim - A Political Simulation Website

**PoliSim** is a Project aiming to deliver a website that can be used with any and all political simulations.
It provides the user and admins, with organisations, titles, documents and other forms of interactive ways to
communicate specifics of the political situation in the simulation.

# Product Information

**PolSim** is build to be as political agnostic as possible. Therefore, it implements general concepts of
political systems and lets the users and admins handle the fine-tuning of their respective simulation. It
does this by organising people into organisation, which can publish Documents/Discussions and Votes. Providing
a Forum for press, handled and reviewable by the administration and a public forum which lets you write short comments
and messages.

## Organisations and Titles

The foundation of any state or even a stateless society is the organisation of people. Organisation are everything
from public forums to private backroom meetings. Almost all forms of organisation can be represented with the developed system
which divides organisation in three groups (public, private and secret). Public organisation have nothing to hide. Their
documents, discussions and votes are public and can be viewed by anyone. Private organisation can make discussions and
votes, when necessary private. Secret organisations can not make anything public. If an organisation changes for example from
private to public all private discussions and votes are made public. This gives a secret organisation, which wants to
announce their existence to the world, the possibility of keeping their past hidden, while posting new public content.

## Press, Letters and Social Media

Besides organisation one of the most important pillars for an immersive simulation is information and contracts.
People can send letters/contracts to other members of the simulation and can write/publish press articles. Either in a
regularly published collection made by the administration or requesting it to be made a breaking news article on its own.
Both ways let the users and administration add context, intrigue and situation which need reactions to the simulation.

It enables Opinion pieces and general discourse which can then be lead in the social media added to the website. As well 
as sparking discussions in the organisations.

# Technical Information

## Building

The project has a public Docker repo where the newest version can be found as a container. An example docker-compose file
can be found in the files of the projects. All needed environment parameters to start up the container are listed here:

```env
DB_NAME=prod
DB_PASSWORD=root
DB_USER=postgres
DB_ADRESS=DB adress with port
ADRESS=host adress with port
INIT_NAME=Root Account Name
INIT_USERNAME=Root Account Username
INIT_PASSWORD=Root Account password
CORS_URL=Your URL needed for CORS
LANG=your language (most likely DE or EN)
```

## Architecture

**PoliSim** is build on 4 layers:
1. Database Layer
2. Data Abstraction Layer
3. Data Validation Layer
4. HTML Serving Layer

**PoliSim** realises the first Layer with GORM. The second layer uses the GORM modells and queries to extract only the needed
data for the request, thereby minimizing the data flow between database and program. The third layer takes in the data from
the HTML-Request, process it, requests the needed complimentary data from the abstraction layer and gives back the state
of validity of the request to the HTML layer to serve an appropriate response.

The HTML requests and responses are received and served with go-chi. The HTML is extended with Hyperscript and HTMX to
minimize the needed HTML send back to the client.

## HTML Serving

HTML is build with the component builder which is inspired by gomponents.
It uses the idea, but simplifies it down to only four functions, which makes shortens the code
quite significantly. It also uses a different approach when validating if the function is an attribute or not.
The HTML snippet builder is seperated from both the validation and serving package.

## Sidebar Handling

Because the website should 

## Language

if you want to support a different language, you have to add a file with your language named 
{your two letter shorthand for your language}.json in the resources folder and change your LANG env variable to that shorthand.
This repo accepts pull requests for any new language support. We as the maintainer will inform you of any new translations needed.