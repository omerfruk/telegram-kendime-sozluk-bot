## Kendime Sözlük Bot

telegram botunu yazdıktan sonra aws ec2 de deploy ettim fakat terminali her kapattıgımda karşılaştıgım sorun malın
durmasi id ben code patlıyor zannediyordum sonradan farkettim ki mal ben her terminali kapattıgımda kendini de
kapatıyormuş bunun çözümü ise bunun service olarak çalışması imiş ben

peki go projesini kendi makinemde nasıl service olarak çalıştırırım

öncelikle projemin build ini aldım ve lokasyonunu kopyaladım

~~~go
 go build main.go
~~~

root yetkisi ile service dosyamı oluşturdum. unutmayalım ki yapılan tüm işlemler 
root yetkisi il oluşturulması lazım yoksa saçma hatalar ile saatlerce uğraşabiliriz


```console
sudo nano /etc/systemd/system/appgo.service
```

oluşan dosyanın içine alttakinleri yapıştırdım

```console
[Unit]
Description=MyApp Go Service
ConditionPathExists=/home/ec2-user/go/pkg/telegram-kendime-sozluk-bot
After=network.target

[Service]
Type=simple
User=root

WorkingDirectory=/home/ec2-user/go/pkg/telegram-kendime-sozluk-bot
ExecStart=/home/ec2-user/go/pkg/telegram-kendime-sozluk-bot/main

Restart=on-failure
RestartSec=10

StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=appgoservice

[Install]
WantedBy=multi-user.target
```

şimdi servis dosyamızı oluşturdugumuza göre bunu service olarak çalıştırmamız lazım 

öncelikle service dosyamızda değişiklik oldugundan dolayı alltaraftaki komutu çalıştırıyoruz 

```console
sudo systemctl daemon-reload
``` 
bu komut yapılan değişiklikleri algılayacaktır

sonrasında sırasi ile  altaraftaki komutları değiştiriyoruz bunlar setvisimizi çalıştırıp çalışma durumunu inceleyecektir 
```console
sudo service appgo start
sudo service appgo status
```

Yeniden başlatmanın ardından otomatik olarak başlatılabilmesi için sistem hizmetine ekleyelim
```console
sudo systemctl enable appgo
sudo systemctl start appgo
```

Böylelikle projemiz service olarak çalışmaya başlıyacaktır 

