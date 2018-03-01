//for replacing underscores in tissue names
String.prototype.replaceAll = String.prototype.replaceAll || function(s, r) {
  return this.replace(new RegExp(s, 'g'), r);
};

//takes an array of values and returns median
function median(values) {

    values.sort( function(a,b) {return a - b;} );
    var half = Math.floor(values.length/2);
    if(values.length % 2)
        return values[half];
    else
        return (values[half-1] + values[half]) / 2.0;
}

// for wrapping text, to get it to work you have to tweak it a bit
function wrap(text, width) {
  text.each(function() {
    var text = d3.select(this),
        words = text.text().split(/\s+/).reverse(),
        word,
        line = [],
        lineNumber = 0,
        lineHeight = 1.1, // ems
        y = text.attr("y"),
        dy = parseFloat(text.attr("dy")),
        tspan = text.text(null).append("tspan").attr("x", 0).attr("y", y).attr("dy", dy + "em");
    while (word = words.pop()) {
      line.push(word);
      tspan.text(line.join(" "));
      if (tspan.node().getComputedTextLength() > width) {
        line.pop();
        tspan.text(line.join(" "));
        line = [word];
        tspan = text.append("tspan").attr("x", 0).attr("y", y).attr("dy", ++lineNumber * lineHeight + dy + "em").text(word);
      }
    }
  });
}

//takes in an array of arrays, and returns an array of means
function mean(values) {
  var average = []
  for (var i = 0; i < values.length; i++) {
    var sum = 0;
    for (var j = 0; j < values[i].length; j++) {
      sum = sum + values[i][j]
    }
    average.push(sum/values[i].length)
  }
  return average;
}

//takes in an array of arrays and an array of means,
// and returns array of std devs
function stdDev(values, means) {
  var stdDevs = []
  for (var i = 0; i < values.length; i++) {
    var temp = 0;
    for (var j = 0; j < values[i].length; j++) {
      temp += Math.pow(values[i][j] - means[i], 2);
    }
    stdDevs.push(Math.sqrt(temp)/values[i].length);
  }
  return stdDevs;
}

//takes in an array of stddevs, and an array of means
//and returns an array of 95% confidence intervals
function confidInterval(stddevs, means) {
  var confids = []
  for (var i = 0; i < stddevs.length; i++) {
    confids.push(means[i] - (1.96*stddevs[i]))
  }
  return confids;
}

//takes in an array of arrays, and returns an array of median absolute deviations
function mad(values) {
  var madResult = [];
  for (var i = 0; i < values.length; i++) {
    var temp = [];
    var med = median(values[i])
    for (var j = 0; j < values[i].length; j++) {
      temp.push(Math.abs(values[i][j] - med))
    }
    madResult.push(median(temp))
  }
  return madResult;
}

//zip arrays together (add more letters for more arrays) to sort
// has an option to cut the array (for waterfall cutting)
// sorts decreasing, to make it increasing:
// sort (function (x, y) {
//   return x.a - y.a
// })
function zip(a, b, c, d, cut) {
  var zipped = [];
  // zip
  for (var i=0; i<a.length; i++) {
    zipped.push({a: a[i], b: b[i], c: c[i], d: d[i] });
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
      c[i] = z.c;
      d[i] = z.d;
    }
  } else { //keep only first and last 20
      var aCopy = [];
      var bCopy = [];
      var cCopy = [];
      var dCopy = [];
      var z;
      for (var i = 0; i < 15; i++) {
        z = zipped[i]
        aCopy.push(z.a);
        bCopy.push(z.b);
        cCopy.push(z.c);
        dCopy.push(z.d);
      }
      //make space for squiggle
      for (var i = 15; i < 20; i++) {
        aCopy.push(0);
        bCopy.push(" ");
        cCopy.push(0);
        dCopy.push([]);
      }
      var count = zipped.length-16
      for (var i = 20; i < 35; i++) {
        z = zipped[count]
        aCopy.push(z.a);
        bCopy.push(z.b);
        cCopy.push(z.c);
        dCopy.push(z.d);
        count++;
      }
      return [aCopy, bCopy, cCopy, dCopy]
  }
}

// makes the AAC waterfall
function makeAAC(svg, drugArrAac, aacMed, aacErr, maxAac, drugIdsAac, width, height, margin, color, plotId, synonyms) {


  //append everything to this: a great trick for visibility
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

      //setting the text format of the x axis
      // go up if negative value, go down if positive value
      // make sure to rotate labels
      // oh yeah and they're clickable, like EVERYTHING ELSE
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
    
    // make a space for middle ticks by removing
    if (aacMed.length > 34) {
      xAxisBar.selectAll(".tick")
      .filter(function(d,i) {
        return ( i < 17) || (i >= 19 ); // leave 17 and 18
      })
        .select("line")
        .remove();

      // make 17 and 18 slanted to represent a cut
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



  //adding chart group for bars
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
            d3.select("#" + drugIdsAac[i] + plotId+ "syn").transition()
            .duration(300).style("opacity", 1);
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function(d,i) {
            d3.select(this).transition()
            .duration(300).style("opacity", 1);
            d3.select("#" + drugIdsAac[i] + plotId + "num").transition()
            .duration(300).style("opacity", 0);
            d3.select("#" + drugIdsAac[i] + plotId + "syn").transition()
            .duration(300).style("opacity", 0);
            d3.select(this).style("cursor", "default");
          }
        })
        .append("svg:title")
          .text(function(d,i) {return drugArrAac[i]});


    //error bars
    var errorBar = chart.selectAll('line.error')
    .data(aacMed);

      errorBar.enter()
        .append('line')
        .attr('class', 'error')
        .style("stroke", "black")
        .attr('x1', function(d,i) { return xrange(i) + (width/aacMed.length)/2; })
        .attr('x2', function(d,i) { return xrange(i) + (width/aacMed.length)/2; })
        .attr('y1', function(d,i) { return yrange(d + aacErr[i]);})
        .attr('y2', function(d,i) { return yrange(d - aacErr[i]);});

      var errorTopBar = chart.selectAll("line.errorTop")
        .data(aacMed);

        errorTopBar.enter()
        .append('line')
        .attr('class', 'errorTop')
        .style("stroke", "black")
        .attr('x1', function(d,i) { if (aacErr[i] != 0) return xrange(i) + (width/aacMed.length)/2-5;})
        .attr('x2', function(d,i) { if (aacErr[i] != 0) return xrange(i) + (width/aacMed.length)/2 + 5; })
        .attr('y1', function(d,i) { if (aacErr[i] != 0) return yrange(d + aacErr[i]);})
        .attr('y2', function(d,i) { if (aacErr[i] != 0) return yrange(d + aacErr[i]);});

        var errorBotBar = chart.selectAll("line.errorBot")
          .data(aacMed);

          errorTopBar.enter()
          .append('line')
          .attr('class', 'errorBot')
          .style("stroke", "black")
          .attr('x1', function(d,i) { if (aacErr[i] != 0) return xrange(i) + (width/aacMed.length)/2-5; })
          .attr('x2', function(d,i) { if (aacErr[i] != 0) return xrange(i) + (width/aacMed.length)/2 + 5; })
          .attr('y1', function(d,i) { if (aacErr[i] != 0) return yrange(d - aacErr[i]);})
          .attr('y2', function(d,i) { if (aacErr[i] != 0) return yrange(d - aacErr[i]);});

    // hover over bar for value
    for (var i = 0; i < aacMed.length; i++) {
      chart.append("text")
        .attr("id", drugIdsAac[i] + plotId + "num")
        .attr("transform", "translate(" + (width) + ",100)")
        .style("text-anchor", "end")
        .style("font-size", "13px")
        .attr("fill", "black")
        .style("opacity", 0)
        .text(drugArrAac[i] + ": " + d3.format(".2f")(aacMed[i]))

        chart.append("text")
        .attr("id", drugIdsAac[i] + plotId + "syn")
        .attr("transform", "translate(" + (width) + ",120)")
        .style("text-anchor", "end")
        .style("font-size", "13px")
        .attr("fill", "black")
        .style("opacity", 0)
        .text("Synonyms: " + synonyms[i].join(", "))
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

// make ic50 waterfall
function makeIC50(svg, drugArrIc50, ic50Med, ic50Err, maxIc50, drugIdsIc50, width, height, margin, color, plotId, synonyms) {

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
          'y':function(d) { return yrange(Math.max(0,d))} // to determine if we should flip the bar or not
        })
        .attr("width", width/ic50Med.length)
        .style('fill', "#2d5faf") // "#2d5faf"

        .on({
          "mouseover": function(d,i) {
            d3.select(this).transition()
            .duration(300).style("opacity", 0.6);
            d3.select("#" + drugIdsIc50[i] + plotId+ "num").transition()
            .duration(300).style("opacity", 1);
            d3.select("#" + drugIdsIc50[i] + plotId+ "syn").transition()
            .duration(300).style("opacity", 1);
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function(d,i) {
            d3.select(this).transition()
            .duration(300).style("opacity", 1);
            d3.select("#" + drugIdsIc50[i] + plotId + "num").transition()
            .duration(300).style("opacity", 0);
            d3.select("#" + drugIdsIc50[i] + plotId + "syn").transition()
            .duration(300).style("opacity", 0);
            d3.select(this).style("cursor", "default");
          }
        })
        .append("svg:title")
          .text(function(d,i) {return drugArrIc50[i]});

          //error bars
          var errorBar = chart.selectAll('line.error')
          .data(ic50Med);

            errorBar.enter()
              .append('line')
              .attr('class', 'error')
              .style("stroke", "black")
            // .merge(errorBar)
              .attr('x1', function(d,i) { return xrange(i) + (width/ic50Med.length)/2; })
              .attr('x2', function(d,i) { return xrange(i) + (width/ic50Med.length)/2; })
              .attr('y1', function(d,i) { return yrange(d + ic50Err[i]);})
              .attr('y2', function(d,i) { return yrange(d - ic50Err[i]);});

              var errorTopBar = chart.selectAll("line.errorTop")
                .data(ic50Med);

                errorTopBar.enter()
                .append('line')
                .attr('class', 'errorTop')
                .style("stroke", "black")
                .attr('x1', function(d,i) { if (ic50Err[i] != 0) return xrange(i) + (width/ic50Med.length)/2-5;})
                .attr('x2', function(d,i) { if (ic50Err[i] != 0) return xrange(i) + (width/ic50Med.length)/2 + 5; })
                .attr('y1', function(d,i) { if (ic50Err[i] != 0) return yrange(d + ic50Err[i]);})
                .attr('y2', function(d,i) { if (ic50Err[i] != 0) return yrange(d + ic50Err[i]);});

                var errorBotBar = chart.selectAll("line.errorBot")
                  .data(ic50Med);

                  errorTopBar.enter()
                  .append('line')
                  .attr('class', 'errorBot')
                  .style("stroke", "black")
                  .attr('x1', function(d,i) { if (ic50Err[i] != 0) return xrange(i) + (width/ic50Med.length)/2-5; })
                  .attr('x2', function(d,i) { if (ic50Err[i] != 0) return xrange(i) + (width/ic50Med.length)/2 + 5; })
                  .attr('y1', function(d,i) { if (ic50Err[i]!= 0) return yrange(d - ic50Err[i]);})
                  .attr('y2', function(d,i) { if (ic50Err[i] != 0) return yrange(d - ic50Err[i]);});

    // hover over bar for value
    for (var i = 0; i < ic50Med.length; i++) {
      chart.append("text")
        .attr("id", drugIdsIc50[i] + plotId + "num")
        .attr("transform", "translate(" + (width) + ",100)")
        .style("text-anchor", "end")
        .style("font-size", "13px")
        .attr("fill", "black")
        .style("opacity", 0)
        .text(drugArrIc50[i] + ": " + d3.format(".2f")(ic50Med[i]))

      chart.append("text")
      .attr("id", drugIdsIc50[i] + plotId + "syn")
      .attr("transform", "translate(" + (width) + ",120)")
      .style("text-anchor", "end")
      .style("font-size", "13px")
      .attr("fill", "black")
      .style("opacity", 0)
      .text("Synonyms: " + synonyms[i].join(", "))
      }

}

// the big function: call this to make both waterfalls!
function makeWaterfall(descriptor, drugArr, aacArr, ic50Arr, plotId, synonyms) {

  var margin = {
    top: 100,
    right: 100,
    bottom: 260,
    left: 80
  };
 
  var synAac = synonyms.slice(0)

  updateDimensions(window.innerWidth, margin);

  // Cut out all NaNs and 0s, and log everything, making a new Ic50 and
  // drug array for ic50

  var temp = [];
  var drugArrIc50 = [];
  var synIc50 = []
  for (var i = 0; i < ic50Arr.length; i++) {
    var temp2 = [];
    for (var j = 0; j < ic50Arr[i].length; j++) {
      if (isNaN(ic50Arr[i][j]) || ic50Arr[i][j] == 0) {

      } else {
        //temp2.push(Math.log10(ic50Arr[i][j]))
        temp2.push((ic50Arr[i][j]))
      }
    }
    if (temp2.length != 0) {
      temp.push(temp2);
      drugArrIc50.push(drugArr[i])
      synIc50.push(synonyms[i])
    }

  }

  ic50Arr = temp;

  //averages
  var aacMed = []
  var aacErr = mad(aacArr)
  /* for confidence intervals */
  //var aacMeans = mean(aacArr);
  //var aacErr = confidInterval(stdDev(aacArr, aacMeans), aacMeans)

  var ic50Med = []
  var ic50Err = mad(ic50Arr)
  //var ic50Means = mean(ic50Arr);
  //var ic50Err = confidInterval(stdDev(ic50Arr, ic50Means), ic50Means)

  for (var i = 0; i < ic50Arr.length; i++) {
    ic50Med.push(median(ic50Arr[i]))
  }
  for (var i = 0; i < aacArr.length; i++) {
    aacMed.push(median(aacArr[i]))
  }
  for (var i = 0; i < aacErr.length; i++) {
    if (aacErr[i] == aacMed[i]) {
      aacErr[i] = 0;
    }
    if (ic50Err[i] == ic50Med[i]) {
      ic50Err[i] = 0;
    }
  }
  var width = 420;
  var height = 500;
  var color = "";

  //zip arrays together so drug arrays are sorted alongside aac/ic50
  //sort increasing, keeping only first and last 20
  var drugArrAacCut, drugArrIc50Cut, synAacCut, synIc50Cut, aacMedCut, aacErrCut, ic50MedCut, ic50ErrCut, maxAac, maxIc50, drugIdsAac, drugIdsIc50, drugArrAacNoCut, drugArrIc50NoCut, maxAacNoCut, maxIc50NoCut, drugIdsAacNoCut, drugIdsIc50NoCut;

  if (drugArr.length > 40) {
    drugArrAacCut = drugArr.slice(0);
    drugArrIc50Cut = drugArrIc50.slice(0);
    synAacCut = synAac.slice(0);
    synIc50Cut = synIc50.slice(0);
    aacMedCut = aacMed.slice(0);
    ic50MedCut = ic50Med.slice(0);
    aacErrCut = aacErr.slice(0);
    ic50ErrCut = ic50Err.slice(0);

    aacTemp = zip(aacMedCut, drugArrAacCut, aacErrCut, synAacCut, true);
    ic50Temp = zip(ic50MedCut, drugArrIc50Cut,  ic50ErrCut, synIc50Cut,true);
    aacMedCut = aacTemp[0]
    ic50MedCut = ic50Temp[0]
    drugArrAacCut = aacTemp[1]
    drugArrIc50Cut = ic50Temp[1]
    aacErrCut = aacTemp[2]
    ic50ErrCut = ic50Temp[2]
    synAacCut = aacTemp[3]
    synIc50Cut = ic50Temp[3]
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
    synAacNoCut = synAac.slice(0);
    synIc50NoCut = synIc50.slice(0);
    zip(aacMed, drugArrAacNoCut, aacErr, synAacNoCut, false);
    zip(ic50Med, drugArrIc50NoCut, ic50Err, synIc50NoCut, false);

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

  // make waterfalls with or without cut
  if (drugArr.length > 40) {
    makeAAC(svg, drugArrAacCut, aacMedCut, aacErrCut, maxAac, drugIdsAac, width, height, margin, color, plotId, synAacCut)
    makeIC50(svg, drugArrIc50Cut, ic50MedCut, ic50ErrCut, maxIc50, drugIdsIc50, width, height, margin, color, plotId, synIc50Cut)
  } else {
    makeAAC(svg, drugArrAacNoCut, aacMed, aacErr, maxAacNoCut, drugIdsAacNoCut, width, height, margin, color, plotId, synAacNoCut)
    makeIC50(svg, drugArrIc50NoCut, ic50Med, ic50Err, maxIc50NoCut, drugIdsIc50NoCut, width, height, margin, color, plotId, synIc50NoCut)
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

    // notes
    if (drugArr.length > 40) {
      svg.append("text")
      .attr("fill", "black")
      .style("text-anchor", "start")
      .attr("dx", -80)
      .attr("dy", height+230)
      .style("font-size", 12)
      .text("* Plot represents the top and bottom 15 data points")
    }

    svg.append("text")
      .attr("fill", "black")
      .style("text-anchor", "start")
      .attr("dx", -80)
      .attr("dy", height+250)
      .style("font-size", 12)
      .text("** Error Bars represent the Median Absolute Deviation")


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
