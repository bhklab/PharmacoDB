function getSupportVec(x, num_points) {
  var dx = Math.pow(10,((Math.log10(d3.max(x)) - Math.log10(d3.min(x)))/(num_points - 1)));
  var support_vec = [];
  for (var i=0; i < num_points; i++){
    support_vec.push((d3.min(x) * Math.pow(dx,i)));
  }
  return support_vec;
}

function hill(x, pars) {
  return (pars[1] + (100 - pars[1]) / (1 + Math.pow(x / pars[2], pars[0])));
}

function ic50(pars) {
  if (50 <= pars[1]) {

    return(NaN)

  }
  // n <- n/100
  //  n <- 1 - n
  // log10(10 ^ pars[3] * ((n - 1) / (pars[2] - n)) ^ (1 / pars[1]))
  return(pars[2] * Math.pow((.5 - 1) / (pars[1]/100 - .5) , (1 / pars[0])));
}

function makeCurveFit(data, pars, minDose, maxDose) {
  //curve fitting data
  var doseArray = data.map(function(x) {return x.dose});
  var x = getSupportVec([minDose, maxDose], 1001);
  var y = [];


  for (var i = 0; i < x.length; i++) {
    y.push(hill(x[i], pars));
  }

  //making into a JSON object to send to curvefit
  var coords = [];
  for (var i = 0; i < x.length; i++) {
    var obj = {}
    obj["x"] = x[i];
    obj["y"] = y[i];
    coords.push(obj);
  }
  return coords;
}

function plotCurveFit (curveCoords, svg, xrange, yrange, width, height, curveId, color, minDose, pars) {
  //set domain and range for curve fit
  var curvefit = d3.svg.line()
    .x(function(d) {return xrange(d.x);})
    .y(function(d) {return yrange(d.y);});

  // add curve fit line
  svg.append("path")
    .attr("id", curveId)
    .attr("d", curvefit(curveCoords))
    .attr("fill", "none")
    // .style("stroke-dasharray", (3,3))
    .attr("stroke", color)
    .attr("stroke-width", 2);

  var area = d3.svg.area()
    .x(function(d) { return xrange(d.x); })
    .y0(0)
    .y1(function(d) { return yrange(d.y); });

  svg.append("path")
    .attr("id", curveId + '-area')
    .datum(curveCoords)
    .attr("class", "area")
    .attr("d", area)
    .attr('fill', color)
    .style('opacity', 0);

  var ec50d = [{x:pars[2], y:hill(pars[2], pars)},
             {x:pars[2], y:0}]

  var ec50d2 = [{x:minDose, y:hill(pars[2], pars)},
             {x:pars[2], y:hill(pars[2], pars)}]

  // console.log(ec50d)

  var ecline = d3.svg.line()
    .x(function(d) {return xrange(d.x);})
    .y(function(d) {return yrange(d.y);})
    .interpolate("linear");

  svg.append("path")
    .attr("class", curveId + '-ec50')
    .attr("d", ecline(ec50d2))
    // .attr("x1", 0)
    // .attr("y1", 0)
    // .attr("x2", 1)
    // .attr("y2", 1)
    .style("stroke", color)
    .style("opacity", 0)
    .attr("stroke-width", 4)
    .style("stroke-dasharray", (3,3));


  var ec50d = [{x:pars[2], y:hill(pars[2], pars)},
             {x:pars[2], y:0}]

  svg.append("path")
    .attr("class", curveId + '-ec50')
    .attr("d", ecline(ec50d))
    // .attr("x1", 0)
    // .attr("y1", 0)
    // .attr("x2", 1)
    // .attr("y2", 1)
    .style("stroke", color)
    .style("opacity", 0)
    .attr("stroke-width", 4)
    .style("stroke-dasharray", (3,3));

  svg.append("circle")
      .attr("class", curveId + '-ec50')
      .attr("r", 6)
      .attr("fill", color)
      .attr("cx", xrange(pars[2]))
      .attr("cy", yrange(hill(pars[2], pars)))
      .style("opacity", 0);

  var ic50d = [{x:ic50(pars), y:hill(ic50(pars), pars)},
           {x:ic50(pars), y:0}]

  var ic50d2 = [{x:minDose, y:50},
           {x:ic50(pars), y:50}]
  // console.log(ic50d)
  // console.log(pars)

  // console.log(ic50d)
  // console.log(pars)
  // console.log(isNaN(ic50(pars)))


  if (!isNaN(ic50(pars))){
    svg.append("path")
    .attr("class", curveId + '-ic50')
    .attr("d", ecline(ic50d))
    // .attr("x1", 0)
    // .attr("y1", 0)
    // .attr("x2", 1)
    // .attr("y2", 1)
    .style("stroke", color)
    .style("opacity", 0)
    .attr("stroke-width", 4)
    .style("stroke-dasharray", (3,3));

    svg.append("path")
    .attr("class", curveId + '-ic50')
    .attr("d", ecline(ic50d2))
    // .attr("x1", 0)
    // .attr("y1", 0)
    // .attr("x2", 1)
    // .attr("y2", 1)
    .style("stroke", color)
    .style("opacity", 0)
    .attr("stroke-width", 4)
    .style("stroke-dasharray", (3,3));

    svg.append("circle")
      .attr("class", curveId + '-ic50')
      .attr("r", 6)
      .attr("fill", color)
      .attr("cx", xrange(ic50(pars)))
      .attr("cy", yrange(50))
      .style("opacity", 0);

  }


}

//DRUG DOSE CURVE
function plotDDRC (data, svg, xrange, yrange, width, height, color, ddrcId, dotsId, trianglesId) {
  // make drug dose response curve
  //does not require data to be bound yet, bound when called in function
  var linepath = d3.svg.line()
    .x(function(d) {return xrange(d.dose);})
    .y(function(d) {return yrange(d.response);});

  // add line
  // svg.append("path")
  //   .attr("id", ddrcId)
  //   .attr("d", linepath(data))
  //   .attr("fill", "none")
  //   .attr("stroke", color)
  //   .style("opacity", 1)
  //   .attr("stroke-width", 2);

  // add plots
  var dots = svg.selectAll("dot")
    .data(data)
    .enter();

    dots.append("circle")
    .filter(function(d) { return d.response <= 100 })
    .filter(function(d) { return 0 <= d.response })
      .attr("id", dotsId)
      .attr("r", 3)
      .attr("fill", color)
      .attr("cx", function(d) {return xrange(d.dose);})
      .attr("cy", function(d) {return yrange(d.response);});

      //eerything above 100, make a triangle and make it at 100 with same xrange
      dots.append("path")
      .filter(function(d) { return d.response > 100 })
        .attr("id", trianglesId)
        .attr("fill", color)
        .attr("d", d3.svg.symbol().type("triangle-up"))
        .attr("transform", function(d) {return "translate(" + xrange(d.dose) + ", " + yrange(100) + ")";});

}

function allIndexesOf (datasets, query) {
  var index = 0, tempIndex = index;
  var arrOfIndexes = []
  while (index != -1) {
    index = datasets.indexOf(query)
    if (index != -1) {
      arrOfIndexes.push(index + tempIndex);
      datasets = datasets.slice(index+1, datasets.length);
      tempIndex = tempIndex + 1 + index
    } else {
      return arrOfIndexes;
    }
  }
  return arrOfIndexes;

}

function updateDimensions(winWidth, margin) {
  width = winWidth - margin.left - margin.right;
  height = 500 - margin.top - margin.bottom;
}


// THE BIG FUNCTION
function makeGraph(data, highlight) {
  //positions and dimensions
  var margin = {
    top: 200,
    right: 300,
    bottom: 100,
    left: 60
  };
  var width = 800;
  var height = 500;
  // var height = 700;
  // updateDimensions(window.innerWidth, margin);


  //calculating min and max of all doses of all datasets
  var minDoseArray = [];
  var maxDoseArray = [];
  var maxResponseArray = [];
  for (var i = 0; i < data.data.length; i++) {
    minDoseArray.push(d3.min(data.data[i].dose_responses, function(d) {return d.dose}));
    maxDoseArray.push(d3.max(data.data[i].dose_responses, function(d) {return d.dose}));
    maxResponseArray.push(d3.max(data.data[i].dose_responses, function(d) {return d.response}));
  }


  var minDose = Math.min.apply(null, minDoseArray);
  var maxDose = Math.max.apply(null, maxDoseArray);
  var maxResponse = Math.max.apply(null, maxResponseArray);

  if (maxResponse < 100) {
    maxResponse = 100;
  }

  //set range for data by domain, and scale by range
  var xrange = d3.scale.log()
    .domain([minDose, maxDose])
    .range([minDose, width]);


  var yrange = d3.scale.linear()
    .domain([0, 100]) // maxResponse
    .range([height, 0]);

  //set axes for graph
  var xAxis = d3.svg.axis()
    .scale(xrange)
    .orient("bottom")
    .tickPadding(2);

  var yAxis = d3.svg.axis()
    .scale(yrange)
    .orient("left")
    .tickPadding(2);

  // Add the svg canvas
  var svg = d3.select("#pmdb_cell_drug")
      .append("svg")
          .attr("width", width + margin.left + margin.right)
          .attr("height", height + margin.top + margin.bottom)
          .attr("id", "ddrcpmdb_cell_drug")
      .append("g")
          .attr("transform",
                "translate(" + margin.left + "," + margin.top + ")")
      .attr("fill", "white");

  // graph title
  svg.append("text")
    .attr("text-anchor", "middle")
    .attr("fill","black")
    .style("font-size", "40px")
    .attr("transform", "translate("+ (width/2) +","+ -100 +")")
    .text(data.data[0].drug.name + " vs. " + data.data[0].cell_line.name);

  // Add the X Axis
  svg.append("g")
      .attr("class", "x axis")
      .attr("transform", "translate(0," + height + ")")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .call(xAxis);

  // X axis label
  svg.append("text")
      .attr("text-anchor", "middle")
      .attr("fill","black")
      .attr("transform", "translate("+ (width/2) +","+(height+40)+")")
      .text("Concentration (uM)");

  // Add the Y Axis
  svg.append("g")
      .attr("class", "y axis")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .call(yAxis);

  // Y axis label
  svg.append("text")
      .attr("text-anchor", "middle")
      .attr("fill","black")
      .attr("transform", "translate("+ -40 +","+(height/2)+")rotate(-90)")
      .text("Response");

  svg.selectAll(".tick").select("text").attr("fill", "black").attr("stroke", "none")

  //for resizing
    // d3.select("svg")
    // .attr( 'preserveAspectRatio',"xMinYMin meet")
    // .attr("viewBox", "0 0 1400 400")
    // .attr('width', 700)
  // .attr('height', height * 8/7 )

  //loop for all instances of data
  // var parsAll = [
  //   [9.387269e-01, 2.051185e+01, 8.787337e-02],
  //   [6.550899e-01, 1.093076e+01, 8.086314e-02]
  // ];

  //organizing the Hill params
  var parsAll = [];
  for (var i = 0; i < data.data.length; i++) {
    parsAll.push([data.data[i].params.HS, data.data[i].params.Einf*100, Math.pow(10,data.data[i].params.EC50)])
  }

  //colors
  var color = d3.scale.category10();
  var similarColor = d3.scale.category20c();

  var datasets = []
  var datasetsForLegend = [];
  var datasetsRaw = [];
  for (var i = 0; i < data.data.length; i++) {
    datasets.push(data.data[i].dataset.name)
    datasetsForLegend.push(data.data[i].dataset.name)
    datasetsRaw.push(data.data[i].dataset.name)

  }

  // doing the rep stuff for the legend label
  for (var i = 0; i < data.data.length; i++) {
    var indsToReplace = allIndexesOf(datasets, datasets[i])
    if (indsToReplace.length > 1) {
      for (var j = 0; j < indsToReplace.length; j++) {
        datasets[indsToReplace[j]] = datasets [indsToReplace[j]] + "rep" + (j+1)
        datasetsForLegend[indsToReplace[j]] = datasetsForLegend[indsToReplace[j]] + " rep " + (j+1) // for legend label
      }
    }
  }


  var nest = d3.nest()
    .key(function(d) {
        return d;
    })
    .entries(datasets);

  var nest2 = d3.nest()
    .key(function(d) {
        return d;
    })
    .entries(datasets);
    console.log(nest)
    var i = 0;
    var j = 0;
    var previous_dataset = "NA";
    //OLD
    // var similarColor = ["#3182bd", "#6baed6", "#9ecae1", "#c6dbef", "#e6550d", "#fd8d3c", "#fdae6b", "#fdd0a2", "#31a354", "#74c476", "#a1d99b", "#c7e9c0", "#756bb1", "#9e9ac8", "#bcbddc", "#dadaeb", "#636363", "#969696", "#bdbdbd", "#d9d9d9"];
    //NEW
    var similarColor =
      [
        '#4b7916', "#6BAC20", "#9EDF53" ,"#C5EC98",
        '#d95f02', "#FD9749", "#FEBA86", "#FEDCC2",
        '#5A54A0', "#7F79B9", "#A6A1CE", "#CDCAE3",
        '#DC187E', "#ED5AA6", "#F17EBA", "#F7B6D7",
        '#0a7bcc', "#1FB0FF", "#85D4FF", "#C2EAFF",
        '#CA9502', "#FDC221", "#FED872", "#FEE9AE",
        '#378169', "#54B697", "#8DCEB9", "#C6E7DC",
        '#525252', "#7A7A7A", "#A3A3A3", "#CCCCCC",
    ]
  while(i < data.data.length) {
  // for (var i = 0; i < data.data.length; i++) {
    //get data for curve fit
    var curveCoords = makeCurveFit(data.data[i].dose_responses, parsAll[i], minDose, maxDose);

    var dotsId = "dots" + datasets[i];
    var trianglesId = "triangles" + datasets[i];
    var curveId = "curve" + datasets[i];

    console.log(datasets[i])
    console.log(previous_dataset)
    // for same datasets, make the color similar hues
    if (datasets[i].indexOf("rep") != -1) {
      // if(previous_dataset == "NA") {
      //   previous_dataset = datasets[i].split("rep")[0]
      // }
      if(datasets[i].split("rep")[0] == previous_dataset){
        j += 1;
        offset = 1;
      } else {
        j = Math.ceil((j+4)/4)*4- 4;
        offset = 0;
      }
      console.log(previous_dataset)
      console.log("j",j)
      console.log("i", i)
      console.log("offset", offset)
      //plot graphs
      plotCurveFit(curveCoords, svg, xrange, yrange, width, height, curveId, similarColor[j - offset], minDose, parsAll[i]);
      plotDDRC (data.data[i].dose_responses, svg, xrange, yrange, width, height, similarColor[j - offset], datasets[i], dotsId);

      //legend colour boxes
      svg.append('rect')
          .attr("x", width+20)
          .attr("y", i*22)
          .attr("width", 20)
          .attr("height", 20)
          .style("fill", similarColor[j - offset]);
          // .style("fill", color(i+1))
      previous_dataset = datasets[i].split("rep")[0]
      i += 1;

    } else {
      j = Math.ceil((j+4)/4)*4 - 4;
      console.log(j)

      plotCurveFit(curveCoords, svg, xrange, yrange, width, height, curveId, similarColor[j], minDose, parsAll[i]);
      plotDDRC(data.data[i].dose_responses, svg, xrange, yrange, width, height, similarColor[j], datasets[i], dotsId, trianglesId);

      //legend colour boxes
      svg.append('rect')
          .attr("x", width+20)
          .attr("y", i*22)
          .attr("width", 20)
          .attr("height", 20)
          .style("fill", function() {
            console.log("fill", similarColor[j])
            return similarColor[j]
          });
          i += 1;
          j += 4;
    }
  }

  // legend labels
  nest.forEach(function(d,i) {
    svg.append('text')
        .attr("x", width+45)
        .attr("y", 16.5 + i * 22)
        .attr("id", "legendLabel" + i)
        .style("text-anchor", "start")
        .style("font-size", 18)
        .attr("fill", "black")
        .on("click", function(){
      		var active   = d.active ? false : true ,
      		  newOpacity = active ? 0 : 1;
      		d3.select("#" + d.key).style("opacity", newOpacity);

          // var newOpacityArea = active ? 0 : 0.7;

          if (active) {
            d3.select("#" + "legendLabel" + i).attr("fill", "silver")
          } else {
            d3.select("#" + "legendLabel" + i).attr("fill", "black")
          }

          d3.selectAll("#" + "dots" + d.key).style("opacity", newOpacity);
          d3.selectAll("#" + "triangles" + d.key).style("opacity", newOpacity);
          d3.selectAll("#" + "curve" + d.key).style("opacity", newOpacity);
          // d3.select("#" + "curve" + d.key + '-area').style("opacity", newOpacityArea);

      		 d.active = active;
        })
        .on({
          "mouseover": function(d) {
            d3.select(this).style("cursor", "pointer"); /*semicolon here*/
          },
          "mouseout": function(d) {
            d3.select(this).style("cursor", "default"); /*and there*/
          }
        })
        .text(datasetsForLegend[i]);
        // console.log(datasetsRaw[i])
        // console.log(highlight.indexOf(datasetsRaw[i]))

        if(highlight.indexOf(datasetsRaw[i])== -1) {
          console.log("Not in highlight")
          var active   = d.active ? false : true ,
            newOpacity = active ? 0 : 1;
          d3.select("#" + d.key).style("opacity", newOpacity);
          // var newOpacityArea = active ? 0 : 0.7;

          if (active) {
            d3.select("#" + "legendLabel" + i).attr("fill", "silver")
          } else {
            d3.select("#" + "legendLabel" + i).attr("fill", "black")
          }

          d3.selectAll("#" + "dots" + d.key).style("opacity", newOpacity);
          d3.selectAll("#" + "triangles" + d.key).style("opacity", newOpacity);
          d3.selectAll("#" + "curve" + d.key).style("opacity", newOpacity);
          // d3.select("#" + "curve" + d.key + '-area').style("opacity", newOpacityArea);

           d.active = active;
        }

        d.aacActive = false;
        d.ec50Active = false;
        d.ic50Active = false;

        d3.select("#" + d.key + "AAC").on({
          "mouseover": function() {
            d3.select("#" + "curve" + d.key + '-area').style("opacity", 0.7); /*semicolon here*/
          },
          "mouseout": function() {
            if(!d.aacActive){
              d3.select("#" + "curve" + d.key + '-area').style("opacity", 0); /*and there*/
            }
          },
            "click": function() {
              newOpacity = d.aacActive ? 0 : 0.7;
              d3.select("#" + "curve" + d.key + '-area').style("opacity", newOpacity);
              d.aacActive = !d.aacActive
            }
        })
        d3.select("#" + d.key + "EC50").on({
          "mouseover": function() {
            d3.selectAll("." + "curve" + d.key + '-ec50').style("opacity", 1); /*semicolon here*/
          },
          "mouseout": function() {
            if(!d.ec50Active){
              d3.selectAll("." + "curve" + d.key + '-ec50').style("opacity", 0); /*and there*/
            }
          },
          "click": function() {
              newOpacity = d.ec50Active ? 0 : 1;
              d3.selectAll("." + "curve" + d.key + '-ec50').style("opacity", newOpacity);
              d.ec50Active = !d.ec50Active
            }
        })
        d3.select("#" + d.key + "IC50").on({
          "mouseover": function() {
            d3.selectAll("." + "curve" + d.key + '-ic50').style("opacity", 1); /*semicolon here*/
          },
          "mouseout": function() {
            if(!d.ic50Active){
              d3.selectAll("." + "curve" + d.key + '-ic50').style("opacity", 0); /*and there*/
            }
          },
          "click": function() {
              newOpacity = d.ic50Active ? 0 : 1;
              d3.selectAll("." + "curve" + d.key + '-ic50').style("opacity", 1);
              d.ic50Active = !d.ic50Active
            }
        })
  });

  nest2.forEach(function(d,i) {
    svg.select("#" + d.key + "AAC")
  })

  var descriptor = data.data[0].drug.name + "-" + data.data[0].cell_line.name;

  d3.select("#pmdb_cell_drug").append("button")
      .attr("type","button")
      .attr("id", "buttonpmdb_cell_drug")
      .attr("class", "downloadButton")
      .text("Download SVG")
      .on("click", function() {
          // download the svg
          downloadSVG("pmdb_cell_drug", descriptor);
      });


}
