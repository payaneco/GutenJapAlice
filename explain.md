# 不思議の英伊のアリスbot

## これは何？

### 概要

これは不思議の国のアリスを段落ごとにイタリア語、英語、日本語でつぶやき続けるBotです。  
イタリア語と英語をだいたい120文字ずつ交互につぶやいた後、最後に日本語をつぶやきます。  
1日4段落くらいつぶやく可能性があります。  

### 詳細

自前のサーバで6時間に1回、1段落分を分割しながらつぶやきますが、動作テストのために手動でつぶやくときもあります。  
1つの段落について、ざっくり10秒間隔でツイートを繰り返します。  
  
ツイートのヘッダは`国旗1-2(3/4)`のように記載されています。  
この意味は「国旗の言語で1章2段落目を4分割したツイートの3個目」です。  
  
機械的に分割したものを表示しているため**各国語訳の段落の位置がずれる可能性があります！**  
そもそも対訳ではありません。  

### ご利用について

このボットは作者個人のイタリア語とGo言語の勉強のためだけに作りました。  
免責を守ったうえで、良識の範囲内で自由にご利用ください。  

## 引用元とか免責とか

### 引用元

#### 英語

英語の底本はProject Gutenbergより"Alice's Adventures in Wonderland by Lewis Carroll"の<a href="http://www.gutenberg.org/files/11/11-0.txt">Plain Text UTF-8</a>を使用させていただいています。  
  
Project Gutenbergのヘッダとフッタを削り、Twitterでの表示のために改行やスペース等の改変を行っています。  
Project Gutenbergと不思議の国のアリスの理念に則り、改変後の文章はパブリックドメイン扱いとなります。  
詳細は上記リンク先のライセンスを参照してください。  
  
#### イタリア語

イタリア語の底本はProject Gutenbergより"Le avventure d'Alice nel paese delle meraviglie by Lewis Carroll"の<a href="http://www.gutenberg.org/cache/epub/28371/pg28371.txt">Plain Text UTF-8</a>を使用させていただいています。  
  
(以下、英語版の表記と同じ)  
Project Gutenbergのヘッダとフッタを削り、Twitterでの表示のために改行やスペース等の改変を行っています。  
Project Gutenbergと不思議の国のアリスの理念に則り、改変後の文章はパブリックドメイン扱いとなります。  
詳細は上記リンク先のライセンスを参照してください。  

#### 日本語

日本語の底本は、プロジェクト杉田玄白よりキャロル、ルイス『不思議の国のアリス』 <a href="http://www.genpaku.org/alice01/alice01j.txt">Text版</a>を使用させていただいています。  

    『不思議の国のアリス』ライセンス
    (C) 1999 山形浩生
    本翻訳は、この版権表示を残す限りにおいて、訳者および著者にたいして許可をとったり使用料を支払ったりすることいっさいなしに、商業利用を含むあらゆる形で自由に利用・複製が認められる。（「この版権表示を残す」んだから、「禁無断複製」とかいうのはダメだぞ）

    プロジェクト杉田玄白　正式参加作品。詳細はhttp://www.genpaku.org/を参照のこと。

電子化や翻訳を行い、無償で公開をしてくださった関係者の皆様の偉業に(日本語でですが)深く感謝いたします。  

### 免責

`@ItenjaBot`を直接的、間接的に使用することにより受けたトラブルや損失・損害等につきましては一切責任を問わないものとします。  
掲載している内容は原本と異なります。  
各国語の段落を関連付ける処理による不整合が随所にみられますが、要望への対応は致しかねます。  
そもそも不思議の国のアリスの妙味である言葉遊びを各国語に訳すとき、対訳せずにローカライズした書き方になっている個所もあります。  
対訳ではありませんので、それを承知の上ご利用ください。  
