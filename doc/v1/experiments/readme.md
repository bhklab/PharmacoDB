# Experiments

```
GET /experiments
```

## Description

This method returns a list of experiments.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- Meta information is included in response headers by default. Use `include` parameter to add info to response body.

## Sources

- http://pharmacodb.pmgenomics.ca/experiments

## Parameters

```
GET /experiments
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **page** | filter | *integer* | no | 1 | The page number for output |
| **per_page** | filter | *integer* | no | 30 | Number of items returned per page |
| **include** | input | metadata | no | - | Include meta info (eg. pagination) in body instead of headers |
| **indent** | input | *boolean* | no | false | Add indentation to response |

## Output Formats

- JSON

## Examples

```
GET /experiments
```

- https://api.pharmacodb.com/v1/experiments

## Output

**JSON**, with metadata included in body.

```
{
    "data": [
        {
            "experiment_id": 1,
            "cell_line": {
                "id": 2,
                "name": "1321N1"
            },
            "tissue": {
                "id": 9,
                "name": "central_nervous_system"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 2,
            "cell_line": {
                "id": 12,
                "name": "22RV1"
            },
            "tissue": {
                "id": 28,
                "name": "prostate"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 3,
            "cell_line": {
                "id": 16,
                "name": "42-MG-BA"
            },
            "tissue": {
                "id": 9,
                "name": "central_nervous_system"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 4,
            "cell_line": {
                "id": 20,
                "name": "5637"
            },
            "tissue": {
                "id": 38,
                "name": "urinary_tract"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 5,
            "cell_line": {
                "id": 24,
                "name": "639-V"
            },
            "tissue": {
                "id": 38,
                "name": "urinary_tract"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 6,
            "cell_line": {
                "id": 26,
                "name": "697"
            },
            "tissue": {
                "id": 13,
                "name": "haematopoietic_and_lymphoid_tissue"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 7,
            "cell_line": {
                "id": 27,
                "name": "769-P"
            },
            "tissue": {
                "id": 15,
                "name": "kidney"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 8,
            "cell_line": {
                "id": 28,
                "name": "786-0"
            },
            "tissue": {
                "id": 15,
                "name": "kidney"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 9,
            "cell_line": {
                "id": 30,
                "name": "8305C"
            },
            "tissue": {
                "id": 36,
                "name": "thyroid"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 10,
            "cell_line": {
                "id": 31,
                "name": "8505C"
            },
            "tissue": {
                "id": 36,
                "name": "thyroid"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 11,
            "cell_line": {
                "id": 29,
                "name": "8-MG-BA"
            },
            "tissue": {
                "id": 9,
                "name": "central_nervous_system"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 12,
            "cell_line": {
                "id": 36,
                "name": "A172"
            },
            "tissue": {
                "id": 9,
                "name": "central_nervous_system"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 13,
            "cell_line": {
                "id": 37,
                "name": "A204"
            },
            "tissue": {
                "id": 33,
                "name": "soft_tissue"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 14,
            "cell_line": {
                "id": 38,
                "name": "A2058"
            },
            "tissue": {
                "id": 31,
                "name": "skin"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 15,
            "cell_line": {
                "id": 39,
                "name": "A253"
            },
            "tissue": {
                "id": 31,
                "name": "skin"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 16,
            "cell_line": {
                "id": 40,
                "name": "A2780"
            },
            "tissue": {
                "id": 24,
                "name": "ovary"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 17,
            "cell_line": {
                "id": 42,
                "name": "A375"
            },
            "tissue": {
                "id": 31,
                "name": "skin"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 18,
            "cell_line": {
                "id": 48,
                "name": "A549"
            },
            "tissue": {
                "id": 18,
                "name": "lung"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 19,
            "cell_line": {
                "id": 49,
                "name": "A673"
            },
            "tissue": {
                "id": 5,
                "name": "bone"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 20,
            "cell_line": {
                "id": 54,
                "name": "ACHN"
            },
            "tissue": {
                "id": 15,
                "name": "kidney"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 21,
            "cell_line": {
                "id": 59,
                "name": "ALL-SIL"
            },
            "tissue": {
                "id": 13,
                "name": "haematopoietic_and_lymphoid_tissue"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 22,
            "cell_line": {
                "id": 62,
                "name": "AMO-1"
            },
            "tissue": {
                "id": 13,
                "name": "haematopoietic_and_lymphoid_tissue"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 23,
            "cell_line": {
                "id": 63,
                "name": "AN3-CA"
            },
            "tissue": {
                "id": 11,
                "name": "endometrium"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 24,
            "cell_line": {
                "id": 67,
                "name": "AsPC-1"
            },
            "tissue": {
                "id": 25,
                "name": "pancreas"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 25,
            "cell_line": {
                "id": 70,
                "name": "AU565"
            },
            "tissue": {
                "id": 7,
                "name": "breast"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 26,
            "cell_line": {
                "id": 72,
                "name": "AZ-521"
            },
            "tissue": {
                "id": 34,
                "name": "stomach"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 27,
            "cell_line": {
                "id": 82,
                "name": "BCPAP"
            },
            "tissue": {
                "id": 36,
                "name": "thyroid"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 28,
            "cell_line": {
                "id": 83,
                "name": "BDCM"
            },
            "tissue": {
                "id": 13,
                "name": "haematopoietic_and_lymphoid_tissue"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 29,
            "cell_line": {
                "id": 89,
                "name": "BFTC-909"
            },
            "tissue": {
                "id": 38,
                "name": "urinary_tract"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        },
        {
            "experiment_id": 30,
            "cell_line": {
                "id": 90,
                "name": "BGC-823"
            },
            "tissue": {
                "id": 34,
                "name": "stomach"
            },
            "compound": {
                "id": 21,
                "name": "AEW541"
            },
            "dataset": {
                "id": 1,
                "name": "CCLE"
            }
        }
    ],
    "metadata": {
        "last_page": 21697,
        "page": 1,
        "per_page": 30,
        "total": 650894
    }
}
```
