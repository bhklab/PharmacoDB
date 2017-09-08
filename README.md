# PharmacoDB

[![Build Status](https://travis-ci.org/bhklab/PharmacoDB.svg?branch=master)](https://travis-ci.org/bhklab/PharmacoDB)
[![Build status](https://ci.appveyor.com/api/projects/status/9bkwyiu0vkm66y1t?svg=true)](https://ci.appveyor.com/project/assefamaru/pharmacodb)
[![GoDoc](https://godoc.org/github.com/bhklab/PharmacoDB/api?status.svg)](https://godoc.org/github.com/bhklab/PharmacoDB/api)

- [About](#about)
- [Example Queries](#example-queries)
- [API](#api)
  - [Endpoints](#endpoints)
  - [Running the API Locally](#running-the-api-locally)
- [Web Application](#web-application)
- [Contributing](#contributing)
- [License](#license)

## About

High throughput drug screening technologies have enabled the profiling of hundreds of cancer cell lines to a large variety of small molecules to discover novel and repurposed treatments. Several large studies have been publicly released testing candidate molecules, often with corresponding molecular profiles of the cell lines used for drug screening. These studies have become invaluable resources for the research community, allowing researchers to leverage the collected data to support their own research. However, such pharmacogenomic datasets are disparate and lack of standardization for cell line and drug identifiers, and used heterogeneous data format for the drug sensitivity measurements.

To address these issues, we developed [PharmacoDB](http://pharmacodb.pmgenomics.ca/), a web-application assembling the largest in vitro drug screens in a single database, and allowing users to easily query the union of studies released to date. PharmacoDB allows scientists to search across publicly available datasets to find instances where a drug or cell line of interest has been profiled, and to view and compare the dose-response data for a specific cell line - drug pair from any of the studies included in the database.

Application currently has two main components: [`api`](./api) and [`web-application`](./front-end).

## Example Queries

#### Cell lines? Try typing 'MCF7'

![MCF7](/front-end/ruby-on-rails/app/assets/images/about/cell-line-search.png)

#### Tissues? Try typing 'Breast'

![Breast](/front-end/ruby-on-rails/app/assets/images/about/tissue-search.png)

#### Drugs? Try typing 'Paclitaxel'

![Paclitaxel](/front-end/ruby-on-rails/app/assets/images/about/drugs-search.png)

#### Drug dose-response curves? Try typing 'MCF7 Paclitaxel'

![MCF7 Paclitaxel](/front-end/ruby-on-rails/app/assets/images/about/drug-dose-response-curve-search.png)

[Start searching](http://pharmacodb.pmgenomics.ca/) across pharmacogenomic datasets and do not hesitate to give us feedback on [GitHub](https://github.com/bhklab/pharmacodb/issues).

## API

The PharmacoDB API contains curated data from a wide set of studies, including: [`CCLE`](http://software.broadinstitute.org/software/cprg/?q=node/11),
[`GDSC`](http://www.cancerrxgene.org/), [`CTRP v2`](https://portals.broadinstitute.org/ctrp/) and [more](http://pharmacodb.pmgenomics.ca/datasets). This data is made restfully available over HTTP.

To access the API, all calls are made to the following URL, adding required parameters for specific services.

```
https://api.pharmacodb.com/v1/
```

Returned data is in `json` format.

Source code documentation for the API can be found at [GoDoc.org](https://godoc.org/github.com/bhklab/PharmacoDB/api)

### Endpoints

**Cell lines**

- [**/cell_lines**](./doc/v1/cell_lines/readme.md)
- [**/cell_lines/{id}**](./doc/v1/cell_lines/id.md)
- [**/cell_lines/{id}/drugs**](./doc/v1/cell_lines/compounds.md)

**Tissues**

- [**/tissues**](./doc/template.md)
- [**/tissues/{id}**](./doc/template.md)
- [**/tissues/{id}/cell_lines**](./doc/template.md)
- [**/tissues/{id}/compounds**](./doc/template.md)

**Drugs**

- [**/compounds**](./doc/template.md)
- [**/compounds/{id}**](./doc/template.md)
- [**/compounds/{id}/cell_lines**](./doc/template.md)
- [**/compounds/{id}/tissues**](./doc/template.md)

**Datasets**

- [**/datasets**](./doc/template.md)
- [**/datasets/{id}**](./doc/template.md)
- [**/datasets/{id}/cell_lines**](./doc/template.md)
- [**/datasets/{id}/tissues**](./doc/template.md)
- [**/datasets/{id}/compounds**](./doc/template.md)

**Experiments**

- [**/experiments**](./doc/template.md)
- [**/experiments/{id}**](./doc/template.md)

**Intersections**

- [**/intersections**](./doc/template.md)
- [**/intersections/1/{cell_id}/{compound_id}**](./doc/template.md)
- [**/intersections/2/{cell_id}/{dataset_id}**](./doc/template.md)

**Stats**

- [**/stats/tissues/cell_lines**](./doc/template.md)
- [**/stats/datasets/cell_lines**](./doc/template.md)
- [**/stats/datasets/tissues**](./doc/template.md)
- [**/stats/datasets/compounds**](./doc/template.md)
- [**/stats/datasets/cell_lines/tissues/{tissue_id}**](./doc/template.md)
- [**/stats/datasets/cell_lines/compounds/{compound_id}**](./doc/template.md)
- [**/stats/datasets/tissues/compounds/{compound_id}**](./doc/template.md)
- [**/stats/datasets/compounds/cell_lines/{cell_id}**](./doc/template.md)
- [**/stats/datasets/compounds/tissues/{tissue_id}**](./doc/template.md)
- [**/stats/datasets/experiments**](./doc/template.md)

Most endpoints contain options for further formatting query or output, including options such as: `indent`, `type`, `include`, `page`, `per_page`, `all` and more. Visit each endpoint link above to see a list of all the options that are available to that endpoint.

### Running the API Locally

To setup and run the API locally, simply download and run one of the api builds already available in the directory [dist](dist). There are builds for various platforms, so pick the executable that corresponds to the OS you wish to run it on.

## Web Application

Documentation for the web-application can be found at: [pharmacodb.pmgenomics.ca/docs](http://pharmacodb.pmgenomics.ca/docs).

## Contributing

If you would like to offer some suggestions, or report any problems regarding the API or web-app, simply create a [new issue](https://github.com/bhklab/PharmacoDB/issues/new) and assign it an appropriate label.

## License

This project is under the GPL-3.0 License - see [LICENSE](LICENSE) for details.
