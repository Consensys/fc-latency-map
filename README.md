# fc-latency-map






### Dependencies:

The project uses <https://github.com/keltia/ripe-atlas>.

A required feature is implemented in a pull request <https://github.com/keltia/ripe-atlas/pull/12>.

* Checkout pr project
```shell
git clone git@github.com:nelsonstr/ripe-atlas.git
cd ripe-atlas
git checkout feature/get-measurement-result-struct%2311
```

Update the pr path in [go.mod](manager/go.mod)


### Configuration:

Must have an ".env" file with the configurations to run.

Is available an example of configuration on [.env.example](./manager/.env.example)

