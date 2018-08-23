# Endpoint

```
GET /{endpoint}
```

## Description

This is a general endpoint doc for endpoints that have the same documentation with a single parameter (`indent`).

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Parameters

```
GET /{endpoint}
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **indent** | input | *boolean* | no | false | Add indentation to response |

## Output Formats

- JSON

## Examples

```
GET /{endpoint}
```

- http://api.pharmacodb.ca/v1/intersections/1/895/526
- http://api.pharmacodb.ca/v1/intersections/2/895/1?indent=true
- http://api.pharmacodb.ca/v1/stats/tissues/cell_lines
- http://api.pharmacodb.ca/v1/stats/datasets/cell_lines
- http://api.pharmacodb.ca/v1/stats/datasets/tissues
- http://api.pharmacodb.ca/v1/stats/datasets/compounds
- http://api.pharmacodb.ca/v1/stats/datasets/experiments
