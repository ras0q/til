# asmuth-bloom

## What I Learned

- 中国剰余定理（CRT）を用いたSecret Sharingの実装
- Mignotte
  - 乱数を使わないため、秘密がある程度大きい必要がある ($M_{(r)} > S > M^{(r-1)}$)
- Asmuth-Bloom
  - 特別なmodを使い、秘密はそれより小さい必要がある
  - 乱数を使う

## Tasks

### Test

```sh
go test -v ./...
```

### Test:asmuth_bloom

```sh
go test -v ./asmuth_bloom/...
```

### Test:mignotte

```sh
go test -v ./mignotte/...
```
