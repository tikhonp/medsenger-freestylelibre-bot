# install go
curl -OL https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz

# Install pkl
sudo curl -L -o /usr/local/bin/pkl https://github.com/apple/pkl/releases/download/0.25.2/pkl-linux-amd64
sudo chmod +x /usr/local/bin/pkl
pkl --version

./build.sh

sudo cp freestylelibre_medsenger_bot.conf /etc/supervisor/conf.d/
sudo cp freestylelibre_nginx.conf /etc/nginx/sites-enabled/
sudo supervisorctl update
sudo systemctl restart nginx
sudo certbot --nginx -d libre.ai.medsenger.ru
