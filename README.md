# Poster
Convert a paper to a poster

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

## 重构中...