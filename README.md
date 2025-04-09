# PoliSim - A Political Simulation Website

![Version is 0.14.1](https://img.shields.io/badge/version-0.14.1-blue)
![Supported Languages are DE, EN](https://img.shields.io/badge/languages-DE,_EN-yellow)
![Repository Licences](https://img.shields.io/github/license/LaybeRize/PoliSim)

**PoliSim** is a Project aiming to deliver a website that can be used with any and all political simulations.
It provides the user and admins, with organisations, titles, documents and other forms of interactive ways to
communicate specifics of the political situation in the simulation.

# Product Information

**PolSim** is build to be as political agnostic as possible. Therefore, it implements general concepts of
political systems and lets the users and admins handle the fine-tuning of their respective simulation. It
does this by organising people into organisation, which can publish Documents/Discussions and Votes. Providing
a forum for press, handled and reviewable by the administration and a public forum which lets you write notes.

## Organisations and Titles

The foundation of any state or even a stateless society is the organisation of people. Organisation are everything
from public forums to private backroom meetings. Almost all forms of organisation can be represented with the developed system
which divides organisation in three groups (public, private and secret). Public organisation have nothing to hide. Their
documents, discussions and votes are public and can be viewed by anyone. Private organisation can make discussions and
votes, when necessary private. Secret organisations can not make anything public. If an organisation changes for example from
secret to public all private discussions and votes from the past are kept hidden. This gives a secret organisation, which 
wants to announce their existence to the world, the possibility of keeping their past hidden, while posting new public content.

## Press, Letters and Boards

Besides organisation one of the most important pillars for an immersive simulation is information and contracts.
People can send letters/contracts to other members of the simulation and can write/publish press articles. Either in a
regularly published collection made by the administration or requesting it to be made a breaking news article on its own.
Both ways let the users and administration add context, intrigue and situation which need reactions to the simulation.

It enables Opinion pieces and general discourse which can be held in the public forum of the website. As well
as sparking discussions in the organisations.

## Chats

Users can also create and write in their own chatrooms, where they can communicate in real time.

# Setup Information

The container composition needs the following environment information to run:
`````
DB_ADDRESS=db
ADDRESS=0.0.0.0:[8080]
NAME=[Your Name Here]
USERNAME=[Your Username Here]
PASSWORD=[test]
LOG_LEVEL=DEBUG
POSTGRES_USER=[user]
POSTGRES_DB=[testdb]
POSTGRES_PASSWORD=[password]
PAGE_NAME=[some string]
ICON_PATH=[path]
WELCOME_FILE=[name]
SQL_KEY=[key value]
`````
`LOG_LEVEL` can be omitted and everything in [] can be customized, the rest is either required to stay that way, or should
stay that way if the docker-compose.yaml is used. DB_PASSWORD and the password after neo4j/ must be the same, otherwise 
the server can't connect to the DB. The password must be at least 8 characters by neo4j's demand. If you are using the 
provided compose file, you can actually change NEO4J_AUTH to `${DB_USER}/${DB_PASSWORD}`.