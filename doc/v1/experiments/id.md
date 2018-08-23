# Experiment

```
GET /experiments/{id}
```

## Description

This method returns a single experiment.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- A `404` error is returned if an item is not found.

## Sources

- http://pharmacodb.ca/experiments

## Parameters

```
GET /experiments/{id}
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **indent** | input | *boolean* | no | false | Add indentation to response |

## Output Formats

- JSON

## Examples

```
GET /experiments/{id}
```

- http://api.pharmacodb.ca/v1/experiments/25
- http://api.pharmacodb.ca/v1/experiments/25?indent=true

## Output

**JSON**

```
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
    },
    "dose_responses": [
        {
            "dose": 0.0025,
            "response": 115.3
        },
        {
            "dose": 0.008,
            "response": 105.67
        },
        {
            "dose": 0.025,
            "response": 107.83
        },
        {
            "dose": 0.08,
            "response": 103.45
        },
        {
            "dose": 0.25,
            "response": 88
        },
        {
            "dose": 0.8,
            "response": 94.1
        },
        {
            "dose": 2.53,
            "response": 80
        },
        {
            "dose": 8,
            "response": 57
        }
    ]
}
```
