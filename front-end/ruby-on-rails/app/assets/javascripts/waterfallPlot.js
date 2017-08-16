//for replacing underscores in tissue names
String.prototype.replaceAll = String.prototype.replaceAll || function(s, r) {
  return this.replace(new RegExp(s, 'g'), r);
};

function median(values) {

    values.sort( function(a,b) {return a - b;} );
    var half = Math.floor(values.length/2);
    if(values.length % 2)
        return values[half];
    else
        return (values[half-1] + values[half]) / 2.0;
}


function zip(a, b, cut) {
  var zipped = [];
  // zip
  for (var i=0; i<a.length; i++) {
    zipped.push({a: a[i], b: b[i]});
  }
  zipped.sort(function (x, y){
    return y.a - x.a;
  });
  // unzip
  if (!cut) {
    var z;
    for (i=0; i<zipped.length; i++) {
      z = zipped[i];
      a[i] = z.a;
      b[i] = z.b;
    }
  } else { //keep only first and last 20
      var aCopy = [];
      var bCopy = [];
      var z;
      for (var i = 0; i < 15; i++) {
        z = zipped[i]
        aCopy.push(z.a);
        bCopy.push(z.b);
      }
      //make space for squiggle
      for (var i = 15; i < 20; i++) {
        aCopy.push(0);
        bCopy.push(" ");
      }
      var count = zipped.length-16
      for (var i = 20; i < 35; i++) {
        z = zipped[count]
        aCopy.push(z.a);
        bCopy.push(z.b);
        count++;
      }
      return [aCopy, bCopy]
  }
}

function makeAAC(svg, drugArrAac, aacMed, maxAac, drugIdsAac, width, height, margin, color, plotId) {


  //append everything to this
  var plot = svg.append("g")
    .attr("id", "plotaac")
    .style("opacity", 0)
    .attr("visibility", "hidden");

  //set range for data by domain, and scale by range
  var xrange = d3.scale.linear()
    .domain([0, drugArrAac.length])
    .range([0, width]);




  var yrange = d3.scale.linear()
    .domain([0, maxAac + 0.01])
    .range([height, 0]);

  //set axes for graph
  var xAxis = d3.svg.axis()
    .scale(xrange)
    .orient("bottom")
    .tickFormat(function(d,i){ return drugArrAac[i] })
    .tickValues(d3.range(drugArrAac.length));

  var yAxis = d3.svg.axis()
    .scale(yrange)
    .orient("left")
    .tickSize(5)
    .tickFormat(d3.format("0.2f"));


  // Add the Y Axis
  var yAxisBar = plot.append("g")
      .attr("class", "y axis")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .call(yAxis)

  // Add the X Axis
  var xAxisBar = plot.append("g")
      .attr("class", "x axis")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .attr("transform", "translate(0," +  height + ")")
      .call(xAxis)

      xAxisBar.selectAll("text")
      .style("text-anchor", "end")
      .style("font-size", "11px")
      .attr("dx", "-.9em")
      .attr("dy", 0)
      .attr("y", (width/aacMed.length)/2)
      .attr("alignment-baseline", "middle")
      .attr("transform", "rotate(-90)")
      .attr("fill", "#207cc1")
      .attr("stroke", "none")
      .on("click", function(i){
            document.location.href = "/search?q=" + drugArrAac[i]
        })
        .on({
          "mouseover": function(d) {
            d3.select(this).transition().duration(400).style("opacity", 0.5);
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function(d) {
            d3.select(this).transition().duration(400).style("opacity", 1);
            d3.select(this).style("cursor", "default");
          }
        })
        .append("svg:title")
          .text(function(d,i) {return drugArrAac[i]});

    if (aacMed.length > 34) {
      xAxisBar.selectAll(".tick")
      .filter(function(d,i) {
        return ( i < 17) || (i >= 19 ); // leave 17 and 18
      })
        .select("line")
        .remove();

      // make 17 and 18 slanted
      xAxisBar.selectAll(".tick")
      .filter(function(d,i) {
        return (i == 17) || (i == 18); // make 17 and 18 slanted
      })
        .select("line")
        .attr("x1", 10)
        .attr("y2", 30)
        .attr("transform", "translate(0,-15)")
    } else {
      xAxisBar.selectAll(".tick")
        .select("line")
        .remove();
    }


      plot.selectAll(".tick")
        .select("text")
        .attr("fill", "black")
        .attr("stroke", "none")



  //adding chart
  var chart = plot.append('g')
      .attr("transform", "translate(0,0)")
      .attr('id','bars')


  // adding each bar
  chart.selectAll('.bar')
      .data(aacMed)
      .enter()
      .append("a")
        .attr("xlink:href", function(d,i) {
          return "/search?q=" + drugArrAac[i]
        })
      .append('rect')
        .attr("class", "bar")
        .style("stroke", "none")
        .attr('width', width/aacMed.length)
        .attr({
          'x':function(d,i){ return xrange(i) }, // each i is the number of the dataset
          'y':function(d){ return yrange(d)} // each d is the actual number of drugs (the num)
        })
        .attr('height',function(d){ return height - yrange(d);})
        .style('fill', "#2d5faf")

        .on({
          "mouseover": function(d,i) {
            d3.select(this).transition()
            .duration(300).style("opacity", 0.6);
            d3.select("#" + drugIdsAac[i] + plotId+ "num").transition()
            .duration(300).style("opacity", 1);
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function(d,i) {
            d3.select(this).transition()
            .duration(300).style("opacity", 1);
            d3.select("#" + drugIdsAac[i] + plotId + "num").transition()
            .duration(300).style("opacity", 0);
            d3.select(this).style("cursor", "default");
          }
        })
        .append("svg:title")
          .text(function(d,i) {return drugArrAac[i]});

    // hover bar for value
    for (var i = 0; i < aacMed.length; i++) {
      chart.append("text")
        .attr("id", drugIdsAac[i] + plotId + "num")
        .attr("transform", "translate(" + (width) + ",100)")
        .style("text-anchor", "end")
        .style("font-size", "13px")
        .attr("fill", "black")
        .style("opacity", 0)
        .text(drugArrAac[i] + ": " + d3.format(".2f")(aacMed[i]))
      }

    d3.select("#plotaac").append("button")
        .attr("type","button")
        .attr("id", "buttonAAC" + plotId)
        .attr("class", "downloadButton")
        .text("Download SVG")
        .on("click", function() {
            // download the svg
            downloadSVG();
        });

}

function makeIC50(svg, drugArrIc50, ic50Med, maxIc50, drugIdsIc50, width, height, margin, color, plotId) {


  //append everything to this
  var plot = svg.append("g")
    .attr("id", "plotic50")
    .style("opacity", 0)
    .attr("visibility", "hidden");

  //set range for data by domain, and scale by range
  var xrange = d3.scale.linear()
    .domain([0, drugArrIc50.length])
    .range([0, width]);

  var yrange = d3.scale.linear()
    .domain([d3.min(ic50Med), maxIc50 + 0.1])
    .range([height, 0]);

  //set axes for graph
  var xAxis = d3.svg.axis()
    .scale(xrange)
    .orient("bottom")
    .tickFormat(function(d,i){ return drugArrIc50[i] })
    .tickValues(d3.range(drugArrIc50.length));


  var yAxis = d3.svg.axis()
    .scale(yrange)
    .orient("left")
    .tickFormat(d3.format("0.2f"));


  // Add the Y Axis
  var yAxisBar = plot.append("g")
      .attr("class", "y axis")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .attr("transform", function() {
        if (d3.min(ic50Med) > 0) {
          return "translate(0," + (yrange(0)-height) + ")"
        }})
      .call(yAxis)



  // Add the X Axis
  var xAxisBar = plot.append("g")
      .attr("class", "x axis")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .attr("transform", function() {
        if (d3.min(ic50Med) > 0) {
          return "translate(0," + yrange(0) + ")"
        } else {
          return "translate(0," +  yrange(0) + ")"
        }
      })
      .call(xAxis)

      xAxisBar.selectAll("text")
      .style("text-anchor", function(d,i) {
        if (ic50Med[i] >= 0) {
          return "end"
        } else {
          return "start"
        }
      })
      .style("font-size", "11px")
      .attr("dx", function(d,i) {
        if (ic50Med[i] >= 0) {
          return "-1em"
        } else {
          return "1em"
        }
      })
      .attr("transform", "rotate(-90)" )
      .attr("dy", 5)
      .attr("y", (width/ic50Med.length)/2)
      .attr("fill", "#207cc1")
      .attr("stroke", "none")
      .on("click", function(d,i){
            document.location.href = "/search?q=" + drugArrIc50[i]
        })
        .on({
          "mouseover": function(d) {
            d3.select(this).transition().duration(400).style("opacity", 0.5);
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function(d) {
            d3.select(this).transition().duration(400).style("opacity", 1);
            d3.select(this).style("cursor", "default");
          }
        })
        .append("svg:title")
          .text(function(d,i) {return drugArrIc50[i]});

      xAxisBar.selectAll(".tick")
      .filter(function(d, i) { return ic50Med[i] < 0; })
        .select("line")
        .attr("transform", "translate(0, -6)");

        if (ic50Med.length > 34) {
          xAxisBar.selectAll(".tick")
          .filter(function(d,i) {
            return ( i < 17) || (i >= 19 ); // leave 17 and 18
          })
            .select("line")
            .remove();

          // make 17 and 18 slanted
          xAxisBar.selectAll(".tick")
          .filter(function(d,i) {
            return (i == 17) || (i == 18); // make 17 and 18 slanted
          })
            .select("line")
            .attr("x1", 10)
            .attr("y2", 30)
            .attr("transform", "translate(0,-15)")
        } else {
          xAxisBar.selectAll(".tick")
            .select("line")
            .remove();
        }



      // xAxisBar.selectAll("text").remove();
      // xAxisBar.selectAll(".tick line").remove();

  plot.selectAll(".tick")
    .select("text")
    .attr("fill", "black")
    .attr("stroke", "none")


  // // for resizing
  // d3.selectAll("#waterfall" + plotId)
  // .attr( 'preserveAspectRatio',"xMinYMin meet")
  // .attr("viewBox", "80 0 700 500")
  // .attr('width', '100%')
  // .attr("height", height )

  //adding chart
  var chart = plot.append('g')
      .attr("transform", "translate(0,0)")
      .attr('id','bars')


  // adding each bar
  chart.selectAll('.bar')
      .data(ic50Med)
      .enter()
      .append("a")
        .attr("xlink:href", function(d,i) {
          return "/search?q=" + drugArrIc50[i]
        })
      .append('rect')
        .attr("class", "bar")
        .style("stroke", "none")
        .attr('height', function(d) {return Math.abs(yrange(d) - yrange(0));})
        .attr({
          'x':function(d,i) { return xrange(i)},
          'y':function(d) { return yrange(Math.max(0,d))}
        })
        .attr("width", width/ic50Med.length)
        .style('fill', "#2d5faf") // "#2d5faf"

        .on({
          "mouseover": function(d,i) {
            d3.select(this).transition()
            .duration(300).style("opacity", 0.6);
            d3.select("#" + drugIdsIc50[i] + plotId+ "num").transition()
            .duration(300).style("opacity", 1);
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function(d,i) {
            d3.select(this).transition()
            .duration(300).style("opacity", 1);
            d3.select("#" + drugIdsIc50[i] + plotId + "num").transition()
            .duration(300).style("opacity", 0);
            d3.select(this).style("cursor", "default");
          }
        })
        .append("svg:title")
          .text(function(d,i) {return drugArrIc50[i]});

    // hover bar for value
    for (var i = 0; i < ic50Med.length; i++) {
      chart.append("text")
        .attr("id", drugIdsIc50[i] + plotId + "num")
        .attr("transform", "translate(" + (width) + ",100)")
        .style("text-anchor", "end")
        .style("font-size", "13px")
        .attr("fill", "black")
        .style("opacity", 0)
        .text(drugArrIc50[i] + ": " + d3.format(".2f")(ic50Med[i]))
      }

}

function makeWaterfall(descriptor, drugArr, aacArr, ic50Arr, plotId) {

  var margin = {
    top: 100,
    right: 100,
    bottom: 200,
    left: 80
  };

  updateDimensions(window.innerWidth, margin);

  // Cut out all NaNs and 0s, and log everything, making a new Ic50 and
  // drug array for ic50
  var temp = [];
  var drugArrIc50 = [];
  for (var i = 0; i < ic50Arr.length; i++) {
    var temp2 = [];
    for (var j = 0; j < ic50Arr[i].length; j++) {
      if (isNaN(ic50Arr[i][j]) || ic50Arr[i][j] == 0) {

      } else {
        temp2.push(Math.log10(ic50Arr[i][j]))
      }
    }
    if (temp2.length != 0) {
      temp.push(temp2);
      drugArrIc50.push(drugArr[i])
    }

  }

  ic50Arr = temp;

  //averages
  var aacMed = []
  var ic50Med = []
  for (var i = 0; i < ic50Arr.length; i++) {
    ic50Med.push(median(ic50Arr[i]))
  }
  for (var i = 0; i < aacArr.length; i++) {
    aacMed.push(median(aacArr[i]))
  }

  var width = 420;
  var height = 500;
  var color = "";

  //zip arrays together so drug arrays are sorted alongside aac/ic50

  //sort increasing, keeping only first and last 20
  var drugArrAacCut, drugArrIc50Cut, aacMedCut, ic50MedCut, maxAac, maxIc50, drugIdsAac, drugIdsIc50, drugArrAacNoCut, drugArrIc50NoCut, maxAacNoCut, maxIc50NoCut, drugIdsAacNoCut, drugIdsIc50NoCut;

  if (drugArr.length > 40) {
    drugArrAacCut = drugArr.slice(0);
    drugArrIc50Cut = drugArrIc50.slice(0);
    aacMedCut = aacMed.slice(0);
    ic50MedCut = ic50Med.slice(0);

    aacTemp = zip(aacMedCut, drugArrAacCut, true);
    ic50Temp = zip(ic50MedCut, drugArrIc50Cut, true);
    aacMedCut = aacTemp[0]
    ic50MedCut = ic50Temp[0]
    drugArrAacCut = aacTemp[1]
    drugArrIc50Cut = ic50Temp[1]

    maxAac = d3.max(aacMedCut);
    maxIc50 = d3.max(ic50MedCut);

    //remove all special characters for the ids of the drugs
    drugIdsAac = [], drugIdsIc50 = []
    for (var i = 0; i < drugArrIc50Cut.length; i++) {
      drugIdsAac.push(drugArrAacCut[i].replaceAll(/[^\w\s]/gi, '').replaceAll(" ", ''))
      drugIdsIc50.push(drugArrIc50Cut[i].replaceAll(/[^\w\s]/gi, '').replaceAll(" ", ''))
    }
  }
  else {

    //sort increasing, without cut
    drugArrAacNoCut = drugArr.slice(0);
    drugArrIc50NoCut = drugArrIc50.slice(0);
    zip(aacMed, drugArrAacNoCut, false);
    zip(ic50Med, drugArrIc50NoCut, false);

    maxAacNoCut = d3.max(aacMed);
    maxIc50NoCut = d3.max(ic50Med);

    //remove all special characters for the ids of the drugs
    drugIdsAacNoCut = [], drugIdsIc50NoCut = []
    for (var i = 0; i < drugArrIc50NoCut.length; i++) {
      drugIdsAacNoCut.push(drugArrAacNoCut[i].replaceAll(/[^\w\s]/gi, '').replaceAll(" ", ''))
      drugIdsIc50NoCut.push(drugArrIc50NoCut[i].replaceAll(/[^\w\s]/gi, '').replaceAll(" ", ''))
    }
  }





  // Add the svg canvas
  var svg = d3.select("#" + plotId)
      .append("svg")
          .attr("width", width + margin.left + margin.right)
          .attr("height", height + margin.top + margin.bottom)
      .attr("id", "waterfall" + plotId)
      .append("g")
          .attr("transform",
                "translate(" + margin.left + "," + margin.top + ")")
      .attr("fill", "white");

  updateDimensions(window.innerWidth, margin);

  // graph title
  var graphTitle = svg.append("text")
    .attr("text-anchor", "middle")
    .attr("fill","black")
    .attr("dy", 0)
    .style("font-size", "23px")
    .attr("transform", "translate("+ (width/2.2) +","+ -50 +")")
    .text(descriptor + ": AAC")
    .call(wrap, width)

  // Y axis label
  var yAxisLabel = svg.append("text")
      .attr("text-anchor", "middle")
      .attr("fill","black")
      .attr("transform", "translate("+ -60 +","+(height/2)+")rotate(-90)")
      .text("AAC");


  if (drugArr.length > 40) {
    makeAAC(svg, drugArrAacCut, aacMedCut, maxAac, drugIdsAac, width, height, margin, color, plotId)
    makeIC50(svg, drugArrIc50Cut, ic50MedCut, maxIc50, drugIdsIc50, width, height, margin, color, plotId)
  } else {
    console.log("ya")
    makeAAC(svg, drugArrAacNoCut, aacMed, maxAacNoCut, drugIdsAacNoCut, width, height, margin, color, plotId)
    makeIC50(svg, drugArrIc50NoCut, ic50Med, maxIc50NoCut, drugIdsIc50NoCut, width, height, margin, color, plotId)
  }


  //determining the view of different plots
  var plotLabels = ["AAC", "pIC50"]
  var plotIds = ["aac", "ic50"]


  //nesting so that each element of Show view has an active property
  var nest = d3.nest()
    .key(function(d) {
        return d;
    })
    .entries([0,1]);

  nest.forEach(function(d,i) {
    var selection = svg.append("text")
      .attr("x", width-50)
      .attr("y", 25 + i * 30)
      .attr("id", "label" + plotIds[i])
      .attr("text-anchor", "start")
      .style("font-size", 18)
      .attr("fill", "black");

    //default
    d3.select("#labelaac").attr("fill", "black");
    d3.select("#labelic50").attr("fill", "silver");
    d3.select("#plotaac").style("opacity", 1);
    d3.select("#plotaac").attr("visibility", "visible");


    //on click of Show views
    selection.on("click", function(){
        var active   = active ? false : true;
        if (active) {
          graphTitle.text(descriptor + ": " + plotLabels[d.key]);
          yAxisLabel.text(plotLabels[d.key]);
          if (d.key == 0) {

          } else {


          }
          d3.select("#" + "label" + plotIds[d.key]).transition().duration(300).attr("fill", "black");
          d3.select("#plot" + plotIds[d.key]).attr("visibility", "visible");
          d3.select("#plot" + plotIds[d.key]).transition().duration(300).style("opacity", 1);

          for (var j = 0; j < 2; j++) {
            if (plotIds[d.key] != plotIds[j]) {
              d3.select("#" + "label" + plotIds[j]).transition().duration(300).attr("fill", "silver");
              d3.select("#plot" + plotIds[j]).attr("visibility", "hidden");
              d3.select("#plot" + plotIds[j]).transition().duration(300).style("opacity", 0);
            }
          }
        } else {
          graphTitle.text();
          yAxisLabel.text();
          d3.select("#" + "label" + plotIds[d.key]).transition().duration(300).attr("fill", "silver");
          d3.select("#plot" + plotIds[d.key]).attr("visibility", "hidden");
          d3.select("#plot" + plotIds[d.key]).transition().duration(300).style("opacity", 0);
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
      .text(plotLabels[i]);
    });



    d3.select("#" + plotId).append("button")
        .attr("type","button")
        .attr("id", "buttonIC50" + plotId)
        .attr("class", "downloadButton")
        .text("Download SVG")
        .on("click", function() {
            // download the svg
            downloadSVG(plotId, descriptor);
        });


}
