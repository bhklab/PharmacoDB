# Compounds

```
GET /cell_lines/{id}/compounds
```

## Description

This method returns a list of unique compounds that have been tested against a cell line of interest.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- Meta information is included in response headers by default. Use `include` parameter to add info to response body.

## Sources

- http://pharmacodb.pmgenomics.ca/cell_lines/

## Parameters

```
GET /cell_lines/{id}/compounds
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **page** | filter | *integer* | no | 1 | The page number for output |
| **per_page** | filter | *integer* | no | 30 | Number of items returned per page |
| **include** | input | metadata | no | - | Include meta info (eg. pagination) in body instead of headers |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | cell_id | Define whether `id = cell_id` or `id = cell_name` |

## Output Formats

- JSON

## Examples

```
GET /cell_lines/{id}/compounds
```

- https://api.pharmacodb.com/v1/cell_lines/895/compounds
- https://api.pharmacodb.com/v1/cell_lines/895/compounds?page=2&per_page=10
- https://api.pharmacodb.com/v1/cell_lines/mcf7/compounds?type=name

## Output

**JSON**, using the cell line `MCF7`, and meta info included in body.

```
{
    "data": [
        {
            "compound": {
                "id": 3,
                "name": "17-AAG"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 7
        },
        {
            "compound": {
                "id": 273,
                "name": "Doxorubicin"
            },
            "datasets": [
                "GDSC1000",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 7
        },
        {
            "compound": {
                "id": 526,
                "name": "paclitaxel"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 7
        },
        {
            "compound": {
                "id": 56,
                "name": "AZD6244"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "compound": {
                "id": 92,
                "name": "Bortezomib"
            },
            "datasets": [
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "compound": {
                "id": 287,
                "name": "Erlotinib"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "compound": {
                "id": 314,
                "name": "Gefitinib"
            },
            "datasets": [
                "GDSC1000",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "compound": {
                "id": 367,
                "name": "Irinotecan"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM"
            ],
            "experiment_count": 6
        },
        {
            "compound": {
                "id": 415,
                "name": "lapatinib"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "compound": {
                "id": 725,
                "name": "Vorinostat"
            },
            "datasets": [
                "GDSC1000",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "compound": {
                "id": 248,
                "name": "Crizotinib"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 271,
                "name": "Docetaxel"
            },
            "datasets": [
                "GDSC1000",
                "gCSI",
                "GRAY"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 313,
                "name": "GDC-0941"
            },
            "datasets": [
                "GDSC1000",
                "gCSI",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 316,
                "name": "Gemcitabine"
            },
            "datasets": [
                "GDSC1000",
                "gCSI",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 328,
                "name": "GSK1070916"
            },
            "datasets": [
                "GDSC1000",
                "GRAY"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 329,
                "name": "GSK1120212"
            },
            "datasets": [
                "GDSC1000",
                "GRAY"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 330,
                "name": "GSK1838705"
            },
            "datasets": [
                "GRAY"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 338,
                "name": "GSK461364"
            },
            "datasets": [
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 438,
                "name": "Methotrexate"
            },
            "datasets": [
                "GDSC1000",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 651,
                "name": "Sorafenib"
            },
            "datasets": [
                "CCLE",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 673,
                "name": "Tamoxifen"
            },
            "datasets": [
                "GDSC1000",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 697,
                "name": "Topotecan"
            },
            "datasets": [
                "CCLE",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "compound": {
                "id": 7,
                "name": "5-FU"
            },
            "datasets": [
                "GDSC1000",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 4
        },
        {
            "compound": {
                "id": 51,
                "name": "AZD-2281"
            },
            "datasets": [
                "GDSC1000",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 4
        },
        {
            "compound": {
                "id": 76,
                "name": "BIBW2992"
            },
            "datasets": [
                "GDSC1000",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 4
        },
        {
            "compound": {
                "id": 197,
                "name": "Carboplatin"
            },
            "datasets": [
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 4
        },
        {
            "compound": {
                "id": 215,
                "name": "CGC-11047"
            },
            "datasets": [
                "GRAY"
            ],
            "experiment_count": 4
        },
        {
            "compound": {
                "id": 237,
                "name": "Cisplatin"
            },
            "datasets": [
                "GDSC1000",
                "GRAY"
            ],
            "experiment_count": 4
        },
        {
            "compound": {
                "id": 290,
                "name": "etoposide"
            },
            "datasets": [
                "GDSC1000",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 4
        },
        {
            "compound": {
                "id": 326,
                "name": "GSK1059615"
            },
            "datasets": [
                "GRAY"
            ],
            "experiment_count": 4
        }
    ],
    "metadata": {
        "last_page": 22,
        "page": 1,
        "per_page": 30,
        "total": 658
    }
}
```
