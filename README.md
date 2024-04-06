<!--suppress HtmlDeprecatedAttribute -->
<div align="center">
    <br>
    <h1>👶 Medsenger CTG monitor bot</h1>
</div>

The __GO__ Medsenger bot for CTG monitors integration.

## 🚀 Install

1. Create DB
2. Create configuration file on `pkl/local/app_config.pkl`
3. RUN `./cmd/setup.sh`

PKL codegen:

# 📦 Development

```bash
pkl-gen-go pkl/config.pkl --base-path github.com/TikhonP/medsenger-freestylelibre-bot
```

Run air:

```bash
air
```

[tailwindcss](https://tailwindcss.com/blog/standalone-cli):

```bash
# Start a watcher
tailwindcss -i view/css/input.css -o public/styles.css --watch

# Compile and minify your CSS for production
tailwindcss -i view/css/input.css -o public/styles.css --minify
```

## 💼 License

Created by Tikhon Petrishchev

Copyright © 2024 OOO Telepat. All rights reserved.

