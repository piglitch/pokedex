# Pokedex CLI

**Pokedex** is an interactive command-line application that allows users to explore the Pokémon world. It enables users to discover locations, catch Pokémon, inspect their details, and manage their collection in a simple and engaging way. 

## Features

- **Explore Locations**: Users can discover various Pokémon locations across the regions.
- **Catch Pokémon**: Interact with Pokémon in your vicinity and add them to your collection.
- **Inspect Pokémon**: View detailed stats and attributes of each Pokémon.
- **Manage Collection**: Keep track of caught Pokémon in your inventory.

## Technical Overview

- **Language**: Go
- **API**: Fetches data from the **PokéAPI** using HTTP requests.
- **Caching**: Implements a custom caching mechanism (`pokecache`) to reduce redundant API calls.
- **CLI Interface**: Allows users to interact with the app via a terminal, ensuring a lightweight experience.
- **Paginated Exploration**: Supports paginated location exploration, making it easy to browse through a large number of items.

This project demonstrates how to use Go for building efficient, user-friendly command-line tools while integrating with external APIs.
