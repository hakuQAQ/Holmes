# Holmes
一款使用cel-go+yaml的指纹识别工具，规则来自Goby <br>
有些yaml规则很不准确，比如：body.icontains("登录")，可自行修改，yaml文件于rule文件夹下。

## Usage
```
Usage: Holmes -v | (-u=<inputurl> | --uf=<inputfile>) (--rf=<rulefile>)... [--th=<threads>] [--to=<timeout>]

FingerPrint Recognition

Options:
  -u, --inputurl   iput a url to scan
      --uf         iput a urlfile to scan
      --th         thread num (default 10)
      --to         Request timeout (default 10)
      --rf         yamlfile path
  -v, --version    Show the version and exit
```

## Build
```
go build .\main.go
```

## Example
```
.\main.exe -u www.baidu.com --rf .\rule\test.yaml
.\main.exe --uf .\1.txt --rf .\rule\
```

![image](https://z3.ax1x.com/2021/11/24/oFAG7Q.png)
