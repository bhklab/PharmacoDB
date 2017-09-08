# Compounds Tested Against A Tissue

```
GET /tissues/{id}/compounds
```

## Description

This method returns a list of unique compounds that have been tested against a tissue of interest.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- Meta information is included in response headers by default. Use `include` parameter to add info to response body.

## Sources

- http://pharmacodb.pmgenomics.ca/tissues

## Parameters

```
GET /tissues/{id}/compounds
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **page** | filter | *integer* | no | 1 | The page number for output |
| **per_page** | filter | *integer* | no | 30 | Number of items returned per page |
| **include** | input | metadata | no | - | Include meta info (eg. pagination) in body instead of headers |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | tissue_id | Define whether `id = tissue_id` or `id = tissue_name` |

## Output Formats

- JSON

## Examples

```
GET /tissues/{id}/compounds
```

- https://api.pharmacodb.com/v1/tissues/7/compounds
- https://api.pharmacodb.com/v1/tissues/7/compounds?page=2&per_page=10
- https://api.pharmacodb.com/v1/tissues/breast/compounds?type=name

## Output

**JSON**, using the tissue **breast**, and meta info included in body.

```
{
    "data": [
        {
            "compound": {
                "id": 526,
                "name": "paclitaxel"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 297
        },
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
            "experiment_count": 296
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
            "experiment_count": 278
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
            "experiment_count": 275
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
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 273
        },
        {
            "compound": {
                "id": 92,
                "name": "Bortezomib"
            },
            "datasets": [
                "GDSC1000",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 270
        },
        {
            "compound": {
                "id": 287,
                "name": "Erlotinib"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 265
        },
        {
            "compound": {
                "id": 330,
                "name": "GSK1838705"
            },
            "datasets": [
                "GRAY"
            ],
            "experiment_count": 265
        },
        {
            "compound": {
                "id": 271,
                "name": "Docetaxel"
            },
            "datasets": [
                "GDSC1000",
                "gCSI",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 252
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
            "experiment_count": 249
        },
        {
            "compound": {
                "id": 329,
                "name": "GSK1120212"
            },
            "datasets": [
                "GDSC1000",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 245
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
            "experiment_count": 244
        },
        {
            "compound": {
                "id": 415,
                "name": "lapatinib"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 243
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
            "experiment_count": 223
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
            "experiment_count": 215
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
            "experiment_count": 212
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
            "experiment_count": 206
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
            "experiment_count": 206
        },
        {
            "compound": {
                "id": 651,
                "name": "Sorafenib"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 204
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
            "experiment_count": 203
        },
        {
            "compound": {
                "id": 507,
                "name": "Nutlin-3"
            },
            "datasets": [
                "CCLE",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 203
        },
        {
            "compound": {
                "id": 76,
                "name": "BIBW2992"
            },
            "datasets": [
                "GDSC1000",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 202
        },
        {
            "compound": {
                "id": 591,
                "name": "Rapamycin"
            },
            "datasets": [
                "GDSC1000",
                "gCSI",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 201
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
            "experiment_count": 199
        },
        {
            "compound": {
                "id": 248,
                "name": "Crizotinib"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 196
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
            "experiment_count": 195
        },
        {
            "compound": {
                "id": 361,
                "name": "Imatinib"
            },
            "datasets": [
                "GDSC1000",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 193
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
            "experiment_count": 190
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
            "experiment_count": 185
        },
        {
            "compound": {
                "id": 523,
                "name": "Oxaliplatin"
            },
            "datasets": [
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 185
        }
    ],
    "metadata": {
        "last_page": 26,
        "page": 1,
        "per_page": 30,
        "total": 759
    }
}
```
