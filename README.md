# RETHi Communication Network Emulator

The RETHi Communication Network Emulator is a tool designed to simulate and test communication networks for the Resilient ExtraTerrestrial Habitats institute (RETHi) project.

[![Project Website](https://img.shields.io/badge/Website-purdue.edu%2Frethi-blue)](https://purdue.edu/rethi)

![Screenshot](./imgs/screenshot.png)

## Features

- Simulate various network conditions
- Test communication protocols in extraterrestrial environments
- Configurable network parameters

## Prerequisites

- Docker
- Docker Compose

## Quick Start

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/rethi-comm-emulator.git
   cd rethi-comm-emulator
   ```

2. Build and run the Docker container:

   ```
   docker build -t amyangxyz111/rethi-comm .
   docker compose up --force-recreate
   ```

3. Access the emulator interface at `http://localhost:8080` (or the appropriate port).

## Configuration

To customize the network settings:

1. Open the `docker-compose.yml` file.
2. Modify the `Local/Remote IP` addresses and ports as needed.
3. Save the file and restart the container.

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Purdue University
- NASA

For more information, visit our [project website](https://purdue.edu/rethi).
