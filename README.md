# chi-go-otoshi

Webサイトのイメージをダウンロードするツール

![](https://6yye61ds.user.webaccel.jp/wp-content/uploads/2016/05/d21f8f61f424bd6f314d96e788b07825.jpg)

*https://yamahack.com/207*


## How to use

create `config.yaml` like below;

```
downloader:
  sites: [
    "gno",
    "anicobin",
  ]
  titles: [
    "大正オトメ御伽話",
    "takt op.Destiny",
    "ジョジョ"
  ]
  saveRootDirectory: "./tmp"
```

run command

```
make build
./downloader
```


