# Poster
It Contains paper summarization and layout.

## Database
We choose the mongodb as the database backend and implemente it. But we use the `Abstract Factory` and `databaser` interface, therefore it's easy to replace the database.

## Users
OAuth2.

## Functions
* Text Summarize    =>  `Text-Rank` Algorithm.
* Question Answering System Handler

## Support
* RESTful API
* Bootstrap Front-end Framework

## Docker
You can use the following command to build a docker image.
```bash
docker build .
```

---
# Below is Chinese Details

## Steps
1. 下载论文
2. 分词，获得词典，包含每个词的IDF值
3. 分句子，用`'.'`、`'。'`分句。
4. 建立PR矩阵，行数为1（句子个数），列数为N。每个元素都为1，或者用一个分布函数，保证和为1。
5. 计算相似度矩阵，用BM25算法，建立一个N*N的对称矩阵，主对角线为0。
6. PR = 0.15 + 0.85 * M * PR[T]
7. 迭代ITER次。

## Tika使用
* 实现go tika client, Docker运行:
```
docker pull logicalspark/docker-tikaserver # only on initial download/update
docker run --rm -p 9998:9998 logicalspark/docker-tikaserver
```

## Posger使用
1. build后运行
2. localhost:8080/generate查看效果

## Bugs:
1. IDF.dict的英文版需要消除大小写的影响
2. 如何获取服务器停止动作，temp: 一个goroutine监控控制台。——需要用来保存IDF.dict
3. 文本摘要如何分段组织。

## 技术难点:
1. pdf 提取图片
2. 把static目录设为根目录， html模版。



