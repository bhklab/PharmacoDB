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
    return (NaN)
  }
  return(pars[2] * Math.pow((.5 - 1) / (pars[1]/100 - .5) , (1 / pars[0])));
}

//Returns the curve-fitting coordinates
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

// plotting the curve fit and all its goodies (AAC, DSS1, EC50, IC50, Einf)
function plotCurveFit (curveCoords, svg, xrange, yrange, width, height, curveId, color, minDose, maxDose, pars) {
  //set domain and range for curve fit
  var curvefit = d3.svg.line()
    .x(function(d) {return xrange(d.x);})
    .y(function(d) {return yrange(d.y);});

  // add curve fit line
  svg.append("path")
    .attr("id", curveId)
    .attr("d", curvefit(curveCoords))
    .attr("fill", "none")
    .attr("stroke", color)
    .attr("stroke-width", 2);

  // AAC Area function
  var area = d3.svg.area()
    .x(function(d) { return xrange(d.x); })
    .y0(0)
    .y1(function(d) { return yrange(d.y); });

  // add AAC area onto graph
  var aac = svg.append("path")
    .attr("id", curveId + '-area')
    .datum(curveCoords)
    .attr("class", "area")
    .attr("d", area)
    .attr('fill', color)
    .style("visibility", "hidden")
    .style('opacity', 0);

  // DSS1 Area function
  var dss1Area = d3.svg.area()
    .x(function(d) {
      return xrange(d.x);
    })
    .y0(function(d) {
      return 50 // means 90
    })
    .y1(function(d) {

      if (yrange(d.y) > 50) {
        return yrange(d.y);
      } else {
        return 50;
      }
    });

  // adding DSS1 area onto graph
  var dss1 = svg.append("path")
    .attr("id", curveId + '-dss1')
    .datum(curveCoords)
    .attr("class", "area")
    .attr("d", dss1Area)
    .attr('fill', color)
    .style("visibility", "hidden")
    .style('opacity', 0);

  // vertical (?) coordinates for EC50 lines
  var ec50d = [{x:pars[2], y:hill(pars[2], pars)},
             {x:pars[2], y:0}]

  // horizontal (?) coordinates for EC50 lines
  var ec50d2 = [{x:minDose, y:hill(pars[2], pars)},
             {x:pars[2], y:hill(pars[2], pars)}]

  // ec50 function
  var ecline = d3.svg.line()
    .x(function(d) {return xrange(d.x);})
    .y(function(d) {return yrange(d.y);})
    .interpolate("linear");

  // adding ec50 horizontal lines onto graph
  svg.append("path")
    .attr("class", curveId + '-ec50a')
    .attr("d", ecline(ec50d2))
    .style("stroke", color)
    .style("opacity", 0)
    .attr("stroke-width", 4)
    .style("stroke-dasharray", (3,3));

  // adding ec50 vertical lines onto graph
  svg.append("path")
    .attr("class", curveId + '-ec50b')
    .attr("d", ecline(ec50d))
    .style("stroke", color)
    .style("opacity", 0)
    .attr("stroke-width", 4)
    .style("stroke-dasharray", (3,3));

  // adding ec50 dot onto graph
  svg.append("circle")
      .attr("class", curveId + '-ec50c')
      .attr("r", 6)
      .attr("fill", color)
      .attr("cx", xrange(pars[2]))
      .attr("cy", yrange(hill(pars[2], pars)))
      .style("opacity", 0);

  // same as the above, but for ic50
  var ic50d = [{x:ic50(pars), y:hill(ic50(pars), pars)},
           {x:ic50(pars), y:0}]

  var ic50d2 = [{x:minDose, y:50},
           {x:ic50(pars), y:50}]

  // if ic50 exists for the curve, add to graph
  if (!isNaN(ic50(pars))){
    svg.append("path")
    .attr("class", curveId + '-ic50a')
    .attr("d", ecline(ic50d))
    .style("stroke", color)
    .style("opacity", 0)
    .attr("stroke-width", 4)
    .style("stroke-dasharray", (3,3));

    svg.append("path")
    .attr("class", curveId + '-ic50b')
    .attr("d", ecline(ic50d2))
    .style("stroke", color)
    .style("opacity", 0)
    .attr("stroke-width", 4)
    .style("stroke-dasharray", (3,3));

    svg.append("circle")
      .attr("class", curveId + '-ic50c')
      .attr("r", 6)
      .attr("fill", color)
      .attr("cx", xrange(ic50(pars)))
      .attr("cy", yrange(50))
      .style("opacity", 0);
}

  // adding Einf onto the graph
  if (!isNaN(pars[1])){
    svg.append("line")
    .attr("class", curveId + '-einf')
    .attr("x1", 0)
    .attr("y1", yrange(pars[1]))
    .attr("x2", xrange(maxDose))
    .attr("y2", yrange(pars[1]))
    .style("stroke", color)
    .style("opacity", 0)
    .attr("stroke-width", 4)
    .style("stroke-dasharray", (3,3));

  }
}

// plotting the dots of the actual drug dose/response curve
function plotDDRC (data, svg, xrange, yrange, width, height, color, ddrcId, dotsId, trianglesId) {
  
  // line function, for when we wanted a line to join the dots
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

  // make group for dots
  var dots = svg.selectAll("dot")
    .data(data)
    .enter();

     // add dots for those that are below 100 for response
     dots.append("circle")
     .filter(function(d) { return d.response <= 100 })
     .filter(function(d) { return 0 <= d.response })
       .attr("id", dotsId)
       .attr("r", 3)
       .attr("fill", color)
       .attr("cx", function(d) {return xrange(d.dose);})
       .attr("cy", function(d) {return yrange(d.response);})
       .on({
         "mouseover": function(d,i) {
           d3.select("." + dotsId + "dose" + i).transition().duration(300).style("opacity", "1");
           d3.select("." + dotsId + "response" + i).transition().duration(300).style("opacity", "1");
         },
         "mouseout": function(d,i) {
           d3.select("." + dotsId + "dose" + i).transition().duration(300).style("opacity", "0");
           d3.select("." + dotsId + "response" + i ).transition().duration(300).style("opacity", "0");
         }
       })

    // add tooltips for the dots: must separate because new lines in d3 are angrifying
    // keep in mind these aren't actual tooltips because making the values appear on the side
    // was easier for me hehe
    // Also more clean imo
    var tooltips = dots.append("g");
      //Dose Tooltip
      tooltips.append("text")
      .filter(function(d) { return d.response <= 100 })
      .filter(function(d) { return 0 <= d.response })
        .attr("class", function(d,i) {return dotsId + "dose" + i})
        .attr("dx", width+20)
        .attr("dy", height/2)
        .attr("font-size", "14px")
        .style("opacity", "0")
          .attr("fill", "black")
          .html(function(d) {return "Dose: " + d3.format(".2f")(d.dose) + " uM"})
      
      //Response tooltip
      tooltips.append("text")
      .filter(function(d) { return d.response <= 100 })
      .filter(function(d) { return 0 <= d.response })
        .attr("class", function(d,i) {return dotsId + "response" + i})
        .attr("dx", width+20)
        .attr("dy", height/2 + 15)
        .attr("font-size", "14px")
        .style("opacity", "0")
          .attr("fill", "black")
          .html(function(d) {return "Response: " + d.response + "%"})


    //everything above 100, make a triangle and make it at 100 with same xrange
    var triangles = svg.selectAll("dot")
    .data(data).enter()
    .append("path")
      .filter(function(d) {
        return d.response > 100
      })
        .attr("id", function(d) {
          return trianglesId;
        })
        .attr("fill", color)
        .attr("d", d3.svg.symbol().type("triangle-up"))
        .attr("transform", function(d) {return "translate(" + xrange(d.dose) + ", " + yrange(100) + ")";})
        .on({
          "mouseover": function(d,i) {
            d3.select("." + trianglesId + "dose" + i).transition().duration(300).style("opacity", "1");
            d3.select("." + trianglesId + "response" + i).transition().duration(300).style("opacity", "1");
          },
          "mouseout": function(d,i) {
            d3.select("." + trianglesId + "dose"+i).transition().duration(300).style("opacity", "0");
            d3.select("." + trianglesId + "response"+i  ).transition().duration(300).style("opacity", "0");
          }
        })

        //triangle tooltips
        var tooltipsTri = dots.append("g");

        tooltipsTri.append("text")
        .filter(function(d) { return d.response > 100 })
          .attr("class", function(d,i) {return trianglesId + "dose" + i})
          .attr("dx", width+20)
          .attr("dy", height/2)
          .attr("font-size", "14px")
          .style("opacity", "0")
            .attr("fill", "black")
            .html(function(d) {return "Dose: " + d3.format(".2f")(d.dose) + " uM"})

        tooltipsTri.append("text")
        .filter(function(d) { return d.response > 100 })
          .attr("class", function(d,i) {return trianglesId + "response" + i})
          .attr("dx", width+20)
          .attr("dy", height/2 + 15)
          .style("opacity", "0")
          .attr("font-size", "14px")
            .attr("fill", "black")
            .html(function(d) {return "Response: " + d.response + "%"})

}

// takes an array of datasets, and finds all indices of where the query occurs
// query being "rep", and returns an array of indices
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

// THE BIG FUNCTION: call this to make a DDRC!
function makeGraph(data, highlight, plotId) {
  
  //positions and dimensions
  var margin = {
    top: 200,
    right: 300,
    bottom: 100,
    left: 60
  };
  var width = 800;
  var height = 500;

  //calculating min and max of all doses of all datasets, and the max response of all datasets
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

  //cut off at 100
  if (maxResponse < 100) {
    maxResponse = 100;
  }

  //set range for data by domain, and scale by range
  var xrange = d3.scale.log()
    .domain([minDose, maxDose])
    .range([minDose, width]);


  var yrange = d3.scale.linear()
    .domain([0, 100])
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
  var svg = d3.select("#" + plotId)
      .append("svg")
          .attr("width", width + margin.left + margin.right)
          .attr("height", height + margin.top + margin.bottom)
          .attr("id", "ddrc" + plotId)
      .append("g")
          .attr("transform",
                "translate(" + margin.left + "," + margin.top + ")");

  // graph title
  var graphTitle = svg.append("text")
    .attr("text-anchor", "middle")
    .attr("id", "ddrcTitle")
    .style("font-size", "40px")
    .attr("transform", "translate("+ (width/2) +","+ -100 +")")
    .html("<a href='/search?q=" + data.data[0].cell_line.name + "' + style='fill:#00bfa5'>" + data.data[0].cell_line.name + "</a>" +  " treated with " + "<a href='/search?q=" + data.data[0].drug.name + "' + style='fill:#00bfa5'>"  + data.data[0].drug.name + "</a>")
    .style("fill","black")
    
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


  //organizing the Hill params
  var parsAll = [];
  for (var i = 0; i < data.data.length; i++) {
    parsAll.push([data.data[i].params.HS, data.data[i].params.Einf*100, Math.pow(10,data.data[i].params.EC50)])
  }

  //colors
  var color = d3.scale.category10();
  var similarColor = d3.scale.category20c();

  // doing the rep stuff for the legend label
  var datasets = []
  var datasetsForLegend = [];
  var datasetsRaw = [];
  for (var i = 0; i < data.data.length; i++) {
    datasets.push(data.data[i].dataset.name)
    datasetsForLegend.push(data.data[i].dataset.name)
    datasetsRaw.push(data.data[i].dataset.name)

  }

  for (var i = 0; i < data.data.length; i++) {
    var indsToReplace = allIndexesOf(datasets, datasets[i])
    if (indsToReplace.length > 1) {
      for (var j = 0; j < indsToReplace.length; j++) {
        datasets[indsToReplace[j]] = datasets [indsToReplace[j]] + "rep" + (j+1)
        datasetsForLegend[indsToReplace[j]] = datasetsForLegend[indsToReplace[j]] + " rep " + (j+1) // for legend label
      }
    }
  }

  //nesting so that we can have objects with attributes (ie .active)
  var nest = d3.nest()
    .key(function(d) {
        return d;
    })
    .entries(datasets);

    var i = 0;
    var j = 0;
    var previous_dataset = "NA";
   
    //a manually curated colour palette wowza
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
    //get data for curve fit
    var curveCoords = makeCurveFit(data.data[i].dose_responses, parsAll[i], minDose, maxDose);

    var dotsId = "dots" + datasets[i];
    var trianglesId = "triangles" + datasets[i];
    var curveId = "curve" + datasets[i];

    // for same datasets, make the color similar hues
    if (datasets[i].indexOf("rep") != -1) {
      if(datasets[i].split("rep")[0] == previous_dataset){
        j += 1;
      } else {
        j = Math.ceil((j+4)/4)*4- 4;
      }
      //plot graphs
      plotCurveFit(curveCoords, svg, xrange, yrange, width, height, curveId, similarColor[j], minDose, maxDose, parsAll[i]); 
      plotDDRC (data.data[i].dose_responses, svg, xrange, yrange, width, height, similarColor[j], datasets[i], dotsId, trianglesId);

      //legend colour boxes
      svg.append('rect')
          .attr("x", width+20)
          .attr("y", i*22)
          .attr("width", 20)
          .attr("height", 20)
          .style("fill", similarColor[j]); 
      
      previous_dataset = datasets[i].split("rep")[0]
      i += 1;

    } else {
      j = Math.ceil((j+4)/4)*4 - 4;

      plotCurveFit(curveCoords, svg, xrange, yrange, width, height, curveId, similarColor[j], minDose, maxDose, parsAll[i]);
      plotDDRC(data.data[i].dose_responses, svg, xrange, yrange, width, height, similarColor[j], datasets[i], dotsId, trianglesId);

      //legend colour boxes
      svg.append('rect')
          .attr("x", width+20)
          .attr("y", i*22)
          .attr("width", 20)
          .attr("height", 20)
          .style("fill", function() {
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

          //to show that it's this dataset has been selected
          if (active) {
            d3.select("#" + "legendLabel" + i).attr("fill", "silver")
          } else {
            d3.select("#" + "legendLabel" + i).attr("fill", "black")
          }

          // hide/show the dataset curve
          d3.selectAll("#" + "dots" + d.key).style("opacity", newOpacity);
          d3.selectAll("#" + "triangles" + d.key).style("opacity", newOpacity);
          d3.selectAll("#" + "curve" + d.key).style("opacity", newOpacity);
          
          // make sure that if any of the summary stats have been shown, hide them all
          if (newOpacity == 0) {
            d3.select("." + "curve" + d.key + '-ec50a').style("opacity", newOpacity);
            d3.select("." + "curve" + d.key + '-ec50b').style("opacity", newOpacity);
            d3.select("." + "curve" + d.key + '-ec50c').style("opacity", newOpacity);
            d3.select("." + "curve" + d.key + '-ic50a').style("opacity", newOpacity);
            d3.select("." + "curve" + d.key + '-ic50b').style("opacity", newOpacity);
            d3.select("." + "curve" + d.key + '-ic50c').style("opacity", newOpacity);
            d3.select("." + "curve" + d.key + '-einf').style("opacity", newOpacity);
            d3.select("#" + "curve" + d.key + '-area').style("opacity", newOpacity);
            d3.select("#" + "curve" + d.key + '-dss1').style("opacity", newOpacity);
          }

      		 d.active = active;
        })
        .on({
          "mouseover": function(d) {
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function(d) {
            d3.select(this).style("cursor", "default"); 
          }
        })
        .text(datasetsForLegend[i]);

        //for when searching [drug] [cell-line] [dataset], where dataset is the highlight
        if(highlight.indexOf(datasetsRaw[i])== -1) {
          console.log("Not in highlight")
          var active   = d.active ? false : true ,
            newOpacity = active ? 0 : 1;
          d3.select("#" + d.key).style("opacity", newOpacity);

          if (active) {
            d3.select("#" + "legendLabel" + i).attr("fill", "silver")
          } else {
            d3.select("#" + "legendLabel" + i).attr("fill", "black")
          }

          d3.selectAll("#" + "dots" + d.key).style("opacity", newOpacity);
          d3.selectAll("#" + "triangles" + d.key).style("opacity", newOpacity);
          d3.selectAll("#" + "curve" + d.key).style("opacity", newOpacity);


           d.active = active;
        }

        d.aacActive = false;
        d.dss1Active = false;
        d.ec50Active = false;
        d.einfActive = false;
        d.ic50Active = false;

        //hovering/clicking effects for everything in the summary stats table
        d3.select("#" + d.key + "AAC").on({
          "mouseover": function() {
            d3.select("#" + "curve" + d.key + '-area').style("opacity", 0.7);
            d3.select("#" + "curve" + d.key + '-area').style("visibility", "visible");
            d3.select("#" + d.key + "AAC.ppl-table-center").transition().duration(300).style("background", "silver")

          },
          "mouseout": function() {
            if(!d.aacActive){
              d3.select("#" + "curve" + d.key + '-area').style("opacity", 0);
              d3.select("#" + "curve" + d.key + '-area').style("visibility", "hidden");
              d3.select("#" + d.key + "AAC.ppl-table-center").transition().duration(300).style("background", "white").style("color", "#6B6B6B")
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
            d3.selectAll("." + "curve" + d.key + '-ec50a').style("opacity", 1); 
            d3.selectAll("." + "curve" + d.key + '-ec50b').style("opacity", 1); 
            d3.selectAll("." + "curve" + d.key + '-ec50c').style("opacity", 1); 
            d3.select("#" + d.key + "EC50.ppl-table-center").transition().duration(300).style("background", "silver")

          },
          "mouseout": function() {
            if(!d.ec50Active){
              d3.selectAll("." + "curve" + d.key + '-ec50a').style("opacity", 0); 
              d3.selectAll("." + "curve" + d.key + '-ec50b').style("opacity", 0); 
              d3.selectAll("." + "curve" + d.key + '-ec50c').style("opacity", 0); 
              d3.select("#" + d.key + "EC50.ppl-table-center").transition().duration(300).style("background", "white").style("color", "#6B6B6B")
            }
          },
          "click": function() {
              newOpacity = d.ec50Active ? 0 : 1;
              d3.selectAll("." + "curve" + d.key + '-ec50a').style("opacity", newOpacity);
              d3.selectAll("." + "curve" + d.key + '-ec50b').style("opacity", newOpacity);
              d3.selectAll("." + "curve" + d.key + '-ec50c').style("opacity", newOpacity);
              d.ec50Active = !d.ec50Active
            }
        })
        d3.select("#" + d.key + "IC50").on({
          "mouseover": function() {
            d3.selectAll("." + "curve" + d.key + '-ic50a').style("opacity", 1); 
            d3.selectAll("." + "curve" + d.key + '-ic50b').style("opacity", 1); 
            d3.selectAll("." + "curve" + d.key + '-ic50c').style("opacity", 1); 
            d3.select("#" + d.key + "IC50.ppl-table-center").transition().duration(300).style("background", "silver")

          },
          "mouseout": function() {
            if(!d.ic50Active){
              d3.selectAll("." + "curve" + d.key + '-ic50a').style("opacity", 0); 
              d3.selectAll("." + "curve" + d.key + '-ic50b').style("opacity", 0); 
              d3.selectAll("." + "curve" + d.key + '-ic50c').style("opacity", 0); 
              d3.select("#" + d.key + "IC50.ppl-table-center").transition().duration(300).style("background", "white").style("color", "#6B6B6B")
            }
          },
          "click": function() {
              newOpacity = d.ic50Active ? 0 : 1;
              d3.selectAll("." + "curve" + d.key + '-ic50a').style("opacity", 1);
              d3.selectAll("." + "curve" + d.key + '-ic50b').style("opacity", 1);
              d3.selectAll("." + "curve" + d.key + '-ic50c').style("opacity", 1);
              d.ic50Active = !d.ic50Active
            }
        })
        d3.select("#" + d.key + "DSS1").on({
          "mouseover": function() {
            d3.select("#" + "curve" + d.key + '-dss1').style("opacity", 0.7);
            d3.select("#" + "curve" + d.key + '-dss1').style("visibility", "visible");
            d3.select("#" + d.key + "DSS1.ppl-table-center").transition().duration(300).style("background", "silver")

          },
          "mouseout": function() {
            if(!d.dss1Active){
              d3.select("#" + "curve" + d.key + '-dss1').style("opacity", 0);
              d3.select("#" + "curve" + d.key + '-dss1').style("visibility", "hidden");
              d3.select("#" + d.key + "DSS1.ppl-table-center").transition().duration(300).style("background", "white").style("color", "#6B6B6B")
            }
          },
            "click": function() {
              newOpacity = d.dss1Active ? 0 : 0.7;
              d3.select("#" + "curve" + d.key + '-dss1').style("opacity", newOpacity);
              d.dss1Active = !d.dss1Active
            }
        })
        d3.select("#" + d.key + "Einf").on({
          "mouseover": function() {
            d3.select("." + "curve" + d.key + '-einf').style("opacity", 0.7);
            d3.select("." + "curve" + d.key + '-einf').style("visibility", "visible");
            d3.select("#" + d.key + "Einf.ppl-table-center").transition().duration(300).style("background", "silver")

          },
          "mouseout": function() {
            if(!d.einfActive){
              d3.select("." + "curve" + d.key + '-einf').style("opacity", 0);
              d3.select("." + "curve" + d.key + '-einf').style("visibility", "hidden");
              d3.select("#" + d.key + "Einf.ppl-table-center").transition().duration(300).style("background", "white").style("color", "#6B6B6B")
            }
          },
            "click": function() {
              newOpacity = d.einfActive ? 0 : 0.7;
              d3.select("." + "curve" + d.key + '-einf').style("opacity", newOpacity);
              d.einfActive = !d.einfActive
            }
        })
  });



  //legend for triangle
  var triangleLeg = svg.append("path")
    .attr("fill", "white")
    .attr("stroke", "black")
    .attr("stroke-width", "1px")
    .attr("d", d3.svg.symbol().type("triangle-up"))
    .attr("transform", function(d) {return "translate(" + (width+30) + ", " + (height-105) + ")";});

    svg.append("text")
      .attr("fill", "black")
      .style("text-anchor", "start")
      .attr("x", width+40)
      .attr("y", height-100)
      .style("font-size", 14)
      .text("=  Values truncated at 100%")

  var descriptor = data.data[0].drug.name + "-" + data.data[0].cell_line.name;

  d3.select("#" +  plotId).append("button")
      .attr("type","button")
      .attr("id", "button" + plotId)
      .attr("class", "downloadButton")
      .text("Download SVG")
      .on("click", function() {
          // download the svg
          $("#ddrcTitle").remove()
          // graph title
          var graphTitle = svg.append("text")
            .attr("text-anchor", "middle")
            .attr("id", "ddrcTitle")
            .style("font-size", "40px")
            .attr("transform", "translate("+ (width/2) +","+ -100 +")")
            .text(data.data[0].cell_line.name + " treated with " + data.data[0].drug.name) 
            .style("fill","black")
          downloadSVG(plotId, descriptor);
      });


}
