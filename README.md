# ZaCherkash

Прокси с хитрыми правилами

# Конфиг

```yaml

# путь до лог файла
log_file: /var/log/zacherkash.log

# слушаем:
bind: 0.0.0.0:8080

# работает также как nginx-upstream:
upstream: 192.168.202.28:9081

# переписываем содержимое body
body_rewrite:
#   rewrite all `google.com` => `google.ru`
    google\.com: google.ru
#   rewrite all `www.opennet.ru` => `ru.opennet.www`
    (www)\.(opennet)\.(ru): $3.$2.$1

```
