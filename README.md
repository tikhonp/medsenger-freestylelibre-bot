<!--suppress HtmlDeprecatedAttribute -->
<div align="center">
    <br>
    <h1>👶 Medsenger CTG monitor bot</h1>
</div>

The __GO__ Medsenger bot for CTG monitors integration.

# 📦 Development

1. Install __docker__ and __make__

2. Create configuration file on `.env`

### Run Development

```sh
make
```

or

```sh
make dev
```

or

```sh
make build-dev # preferred if config files were changed, so it rebuilds image
```

### HTML templating

I use [templ](https://github.com/a-h/templ) as template engine. After changing `*.templ` files regenerate go code using:

```sh
make templ
```

> development docker container must be active

### Enter server container shell

There is shortcut for this:

```sh
make go-to-server-container
```

[tailwindcss](https://tailwindcss.com/blog/standalone-cli):

```bash
# Start a watcher
tailwindcss -i view/css/input.css -o public/styles.css --watch

# Compile and minify your CSS for production
tailwindcss -i view/css/input.css -o public/styles.css --minify
```

# Deploying

To deploy you also need __docker__ and __make__. In project root run:

```sh
make prod
```

It will create prod containers and run it in detached mode.

To stop run:

```sh
make fprod
```

To view logs in real time:

```sh
make logs-prod
```

# 💼 License

Created by Tikhon Petrishchev

Copyright © 2025 OOO Telepat. All rights reserved.

