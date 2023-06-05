RETHi-MCVT-HMS-CN
========================

Purpose
-------

Dockerized communicaiton network subsystem for MCVT v6.3. 

Get started
---------------------

CN only
^^^^^^^^

Clone this repo and run :code:`docker compose up`, it will automatically download the compiled container from docker-hub and run it with default environment variables.

With DRDS, C2, MCVT
^^^^^^^^^^^^^^^^^^^

Run :code:`docker compose up` with the :code:`docker-compose.yml` from the DRDS repo.

Build the container
---------------------

**Local use**: :code:`docker build -t comm .`

**Push to docker-hub** :code:`docker buildx build -t amyangxyz111/rethi-comm --platform linux/arm64,linux/amd64 --push .`

