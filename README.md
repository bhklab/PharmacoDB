# PharmacoDb API

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bhklab/PharmacoDb-api/blob/master/LICENSE)

This is the official PharmacoDb API, enabling development of applications using curated data from multiple datasets.

## Accessing the API

All calls are made to the following URL, adding required parameters/endpoints for specific services.

```
https://api.pharmacodb.com/v1/
```

Returned data is in `json` format.

## Endpoints

### Cell Lines

* /cell_lines
* /cell_lines/stats
* /cell_lines/ids
* /cell_lines/ids/:id
* /cell_lines/ids/:id/drugs
* /cell_lines/ids/:id/drugs/stats
* /cell_lines/names
* /cell_lines/names/:name
* /cell_lines/names/:name/drugs
* /cell_lines/names/:name/drugs/stats

### Tissues

* /tissues
* /tissues/stats
* /tissues/ids
* /tissues/ids/:id
* /tissues/ids/:id/cell_lines
* /tissues/ids/:id/drugs
* /tissues/names
* /tissues/names/:name
* /tissues/names/:name/cell_lines
* /tissues/names/:name/drugs

### Drugs

* /drugs
* /drugs/stats
* /drugs/ids
* /drugs/ids/:id
* /drugs/ids/:id/cell_lines
* /drugs/ids/:id/tissues
* /drugs/names
* /drugs/names/:name
* /drugs/names/:name/cell_lines
* /drugs/names/:name/tissues

### Datasets

* /datasets
* /datasets/stats
* /datasets/ids
* /datasets/ids/:id
* /datasets/ids/:id/cell_lines
* /datasets/ids/:id/tissues
* /datasets/ids/:id/drugs
* /datasets/names
* /datasets/names/:name
* /datasets/names/:name/cell_lines
* /datasets/names/:name/tissues
* /datasets/names/:name/drugs

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) for details.
