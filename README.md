[![Go Report Card](https://goreportcard.com/badge/github.com/kurogee/wagyu-script)](https://goreportcard.com/report/github.com/kurogee/wagyu-script)

![Wagyu-script's logo](./Wagyu-script_logo.png)

Wagyu-script
---
Wagyu-scriptは、@kurogeが開発した半角スペース区切りが特徴のシンプルなプログラミング言語です。/ Wagyu-script is a simple programming language developed by @kuroge, which is characterized by half-width space separation.

始め方など
---
詳しくは、[Get Started](./docs/ja/get-started.md)ページをご覧ください。

How to start (en)
---
Please see the [Get Started](./docs/en/get-started.md) page for more information (Used Google Translate).

**Notice**: The Japanese page may explain how to use it based on the latest version. If you feel that the information in the English version is out of date, please translate the Japanese version of the page.

お知らせ
---
 - 2024/08/06: Wagyu-scriptのパッケージの反映忘れが発覚したため、v1.1.0-betaより前（v1.1.0-betaは**含まない**という意味）のバージョンが非推奨となりました。

その他細かい更新内容は、[CHANGELOG.md](./CHANGELOG.md)をご覧ください。

Examples
---

### Hello, World!
```
println "Hello, World!";
```

### FizzBuzz
```
each #from(1 100) > i {
    if #(i % 15 == 0) {
        println "FizzBuzz";
    } elif #(i % 5 == 0) {
        println "Buzz";
    } elif #(i % 3 == 0) {
        println "Fizz";
    } else {
        println i;
    };
};
```

Credits
---
This project uses codes from the following projects:
 - [github.com/Knetic/govaluate](https://github.com/Knetic/govaluate) (MIT License)
    Some functions are used in this codes. Thanks for the great library! / このコードには以上のライブラリの関数が使用されています。これらのライブラリの製作者に感謝します！

License
---
This project is licensed under the MIT License. Please see the [LICENSE](./LICENSE) file for more information. / このプロジェクトはMITライセンスのもとで公開されています。詳しくは[LICENSE](./LICENSE)ファイルをご覧ください。