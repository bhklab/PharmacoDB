**PharmacoDB** database and web-application. [http://pharmacodb.pmgenomics.ca/](http://pharmacodb.pmgenomics.ca/)

[![Build Status](https://travis-ci.org/bhklab/PharmacoDB.svg?branch=master)](https://travis-ci.org/bhklab/PharmacoDB)
[![Build status](https://ci.appveyor.com/api/projects/status/9bkwyiu0vkm66y1t?svg=true)](https://ci.appveyor.com/project/assefamaru/pharmacodb)
[![GoDoc](https://godoc.org/github.com/bhklab/PharmacoDB/api?status.svg)](https://godoc.org/github.com/bhklab/PharmacoDB/api)
[![Go Report Card](https://goreportcard.com/badge/github.com/bhklab/PharmacoDB)](https://goreportcard.com/report/github.com/bhklab/PharmacoDB)

[PharmacoDB](http://pharmacodb.pmgenomics.ca/) enables efficient mining of a compendium of large pharmacogenomic studies where panels of immortalized cancer cell lines have been screened against hundreds of approved and experimental drug compounds. Database contains curated datasets with standard annotations for cell lines, their tissue source, and drug compounds, as well as hundreds of thousands of drug dose-response curves.

Application currently has two main components: [`api`](./api) and [`web-application`](./front-end).

## Table of Contents

- [**API**](#api)
  - [Endpoints](#endpoints)
  - [Running the API Locally](#running-the-api-locally)
- [**Web Application**](#web-application)
- [**Contributing**](#contributing)
- [**License**](#license)

## API

The PharmacoDB API contains curated data from a wide set of studies, including: [`CCLE`](http://software.broadinstitute.org/software/cprg/?q=node/11),
[`GDSC`](http://www.cancerrxgene.org/), [`CTRP v2`](https://portals.broadinstitute.org/ctrp/) and [more](http://pharmacodb.pmgenomics.ca/datasets). This data is made restfully available over HTTP.

To access the API, all calls are made to the following URL, adding required parameters for specific services.

```
https://api.pharmacodb.com/v1/
```

Returned data is in `json` format.

### Endpoints

Resource type: **`cell lines`**

- [**/cell_lines**](./doc/api/template.md)
- [**/cell_lines/{id}**](./doc/api/template.md)
- [**/cell_lines/{id}/drugs**](./doc/api/template.md)

Resource type: **`tissues`**

- [**/tissues**](./doc/api/template.md)
- [**/tissues/{id}**](./doc/api/template.md)
- [**/tissues/{id}/cell_lines**](./doc/api/template.md)
- [**/tissues/{id}/drugs**](./doc/api/template.md)

Resource type: **`drugs`**

- [**/drugs**](./doc/api/template.md)
- [**/drugs/{id}**](./doc/api/template.md)
- [**/drugs/{id}/cell_lines**](./doc/api/template.md)
- [**/drugs/{id}/tissues**](./doc/api/template.md)

Resource type: **`datasets`**

- [**/datasets**](./doc/api/template.md)
- [**/datasets/{id}**](./doc/api/template.md)
- [**/datasets/{id}/cell_lines**](./doc/api/template.md)
- [**/datasets/{id}/tissues**](./doc/api/template.md)
- [**/datasets/{id}/drugs**](./doc/api/template.md)

Resource type: **`experiments`**

- [**/experiments**](./doc/api/template.md)
- [**/experiments/{id}**](./doc/api/template.md)

Resource type: **`intersections`**

- [**/intersections**](./doc/api/template.md)
- [**/intersections/1/{cell_id}/{drug_id}**](./doc/api/template.md)
- [**/intersections/2/{cell_id}/{dataset_id}**](./doc/api/template.md)

Most endpoints contain options for further formatting query or output, including options such as: `indent`, `type`, `include`, `page`, `per_page`, `all` and more. Visit each endpoint link above to see a list of all the options that are available to that endpoint.

### Running the API Locally

To setup and run the API locally, simply download and run one of the api builds already available in the directory [dist/api](dist/api). There are builds for various platforms, so pick the executable that corresponds to the OS you wish to run it on.

## Web Application

Add webapp content here for documenting some of the interface features, and linking to the web docs page.

## Contributing

If you would like to offer some suggestions, or report any problems regarding the API or web-app, simply create a [new issue](https://github.com/bhklab/PharmacoDB/issues/new) and assign it an appropriate label.

## License

This project is under the MIT License - see [LICENSE](LICENSE) for details.
