function updateDimensions(winWidth, margin) {
  width = winWidth - margin.left - margin.right;
  height = 500 - margin.top - margin.bottom;
}

// takes two arrays of values and returns an array of intersecting values
function findIntersection(set1, set2) {
  //see which set is shorter
  var temp;
  if (set2.length > set1.length) {
      temp = set2, set2 = set1, set1 = temp;
  }

  return set1
    .filter(function(e) { //puts in the intersecting names
      return set2.indexOf(e) > -1;
    })
    .filter(function(e,i,c) { // gets rid of duplicates
      return c.indexOf(e) === i;
    })
}

// returns coordinates for drawing the intersections
function intersection(x0, y0, r0, x1, y1, r1) {
      var a, dx, dy, d, h, rx, ry;
      var x2, y2;

      /* dx and dy are the vertical and horizontal distances between
       * the circle centers.
       */
      dx = x1 - x0;
      dy = y1 - y0;

      /* Determine the straight-line distance between the centers. */
      d = Math.sqrt((dy * dy) + (dx * dx));

      /* Check for solvability. */
      if (d > (r0 + r1)) {
        /* no solution. circles do not intersect. */
        return false;
      }
      if (d < Math.abs(r0 - r1)) {
        /* no solution. one circle is contained in the other */
        return false;
      }

      /* 'point 2' is the point where the line through the circle
       * intersection points crosses the line between the circle
       * centers.
       */

      /* Determine the distance from point 0 to point 2. */
      a = ((r0 * r0) - (r1 * r1) + (d * d)) / (2.0 * d);

      /* Determine the coordinates of point 2. */
      x2 = x0 + (dx * a / d);
      y2 = y0 + (dy * a / d);

      /* Determine the distance from point 2 to either of the
       * intersection points.
       */
      h = Math.sqrt((r0 * r0) - (a * a));

      /* Now determine the offsets of the intersection points from
       * point 2.
       */
      rx = -dy * (h / d);
      ry = dx * (h / d);

      /* Determine the absolute intersection points. */
      var xi = x2 + rx;
      var xi_prime = x2 - rx;
      var yi = y2 + ry;
      var yi_prime = y2 - ry;

      return [xi, xi_prime, yi, yi_prime];

}

//for the difference of arrays - particularly in the intersections and middles
//does not mutate any of the arrays
Array.prototype.diff = function(a) {
    return this.filter(function(i) {return a.indexOf(i) < 0;});
};

//for replacing underscores in tissue names
String.prototype.replaceAll = String.prototype.replaceAll || function(s, r) {
  return this.replace(new RegExp(s, 'g'), r);
};

// creates the shapes of the circles/intersections for the venn diagram
// this entire function is a horror story
// petr/victor/whoever's reading this, just ask me if you need anything in this function
// i'll decipher it
function makeCircles(svg, color, arrColor, numCircles, names, vennId, datasets, middle) {

  //determine circle dimensions
  var rad = 200,
      leftX = 50,
      rightX = 260,
      topY = 180,
      botMid = 150,
      botY = 380;

  //append everything to venn so that it can be opacity 0
  var venn = svg.append("g")
    .attr("id", "venn" + vennId)
    .attr("transform", "translate(-80, 0)")
    .style("opacity", 0)
    .attr("visibility", "hidden");

  //make circles
  var circle0 = venn.append("circle")
      .attr("cx", leftX)
      .attr("cy", topY)
      .attr("r", rad)
      .attr("id", vennId + "circle0")
      .style("opacity", 1)
      .attr("fill", color(0))
      .on({
        "mouseover": function() {
          d3.select("#" + vennId + "circle0").transition()
          .duration(300).style("opacity", 0.5); // slice
          d3.select(this).style("cursor", "pointer");
        },
        "mouseout": function() {
          d3.select("#" + vennId + "circle0").transition()
          .duration(300).style("opacity", 1); // slice opacity
          d3.select(this).style("cursor", "default");
        }
      });


      //circle label
      venn.append("text")
          .attr("transform", "translate(-100,-30)")
          .attr("text-anchor", "middle")
          .attr("fill", "black")
          .style("font-size", "25px")
          .text(datasets[0] + " (" + names[0].length + ")");

  var circle1 = venn.append("circle")
      .attr("cx", rightX)
      .attr("cy", topY)
      .attr("r", rad)
      .attr("id", vennId + "circle1")
      .style("opacity", 1)
      .attr("fill", color(1))
      .on({
        "mouseover": function() {
          d3.select("#" + vennId + "circle1").transition()
          .duration(300).style("opacity", 0.5); // slice
          d3.select(this).style("cursor", "pointer");
        },
        "mouseout": function() {
          d3.select("#" + vennId + "circle1").transition()
          .duration(300).style("opacity", 1); // slice opacity
          d3.select(this).style("cursor", "default");
        }
      });

      //circle label
      venn.append("text")
          .attr("transform", "translate(380,-30)")
          .attr("text-anchor", "middle")
          .attr("fill", "black")
          .style("font-size", "25px")
          .text(datasets[1] + " (" + names[1].length + ")");

  //compute intersection01
  var int01 = findIntersection(names[0], names[1]);
  var int01NoMid = int01.diff(middle);
  var intPt01 = intersection(leftX, topY, rad, rightX, topY, rad);
  var intShape01 = venn.append("g")
      .append("path")
      .attr("id", vennId + "intShape01")
      .attr("d", function() {
        return "M" + intPt01[0] + "," + intPt01[2] + "A" + rad + "," + rad +
          " 0 0,1 " + intPt01[1] + "," + intPt01[3]+ "A" + rad + "," + rad +
          " 0 0,1 " + intPt01[0] + "," + intPt01[2];
      })
      .style("opacity", 1)
      .style('fill', color(2))
      .on("click", function(){
        $("table").remove();
        if (numCircles == 2) {
          makeTable(int01, vennId + "int01")
        } else if (numCircles == 3){
          makeTable(int01NoMid, vennId + "int01")
        }

        d3.select("#" + vennId + "int01").transition().duration(400).style("opacity", 1);
      })
      .on({
        "mouseover": function() {
          d3.select("#" + vennId + "intShape01").transition()
          .duration(300).style("opacity", 0.5); // slice
          d3.select(this).style("cursor", "pointer");
        },
        "mouseout": function() {
          d3.select("#" + vennId + "intShape01").transition()
          .duration(300).style("opacity", 1); // slice opacity
          d3.select(this).style("cursor", "default");
        }
      });

  //intersection label    
  var topMidCount = venn.append("text")
          .attr("transform", "translate(155,170)")
          .attr("text-anchor", "middle")
          .attr("fill", "white")
          .style("font-size", "20px")
          .text(int01.length);
  
  // circle label
  var circle0Count = venn.append("text")
      .attr("transform", "translate(-30,170)")
      .attr("text-anchor", "middle")
      .attr("fill", "white")
      .style("font-size", "20px")
      .text(names[0].length - int01.length);

  // circle label
  var circle1Count = venn.append("text")
      .attr("transform", "translate(330,170)")
      .attr("text-anchor", "middle")
      .attr("fill", "white")
      .style("font-size", "20px")
      .text(names[1].length - int01.length);

  var threeCNoMid0;
  var threeCNoMid1;
  var threeCNoMid2;

  //if three circles
  if (numCircles == 3) {
    var circle2 = venn.append("circle")
        .attr("cx", botMid)
        .attr("cy", botY)
        .attr("r", rad)
        .attr("id", vennId + "circle2")
        .style("opacity", 1)
        .attr("fill", color(3))
        .on({
          "mouseover": function() {
            d3.select("#" + vennId + "circle2").transition()
            .duration(300).style("opacity", 0.5); // slice
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function() {
            d3.select("#" + vennId + "circle2").transition()
            .duration(300).style("opacity", 1); // slice opacity
            d3.select(this).style("cursor", "default");
          }
        });

        //circle label
        venn.append("text")
            .attr("transform", "translate(160,630)")
            .attr("text-anchor", "middle")
            .attr("fill", "black")
            .style("font-size", "25px")
            .text(datasets[2] + " (" + names[2].length + ")");

    //intersection02
    var int02 = findIntersection(names[0], names[2]);
    var int02NoMid = int02.diff(middle);
    var intPt02 = intersection(leftX, topY, rad, botMid, botY, rad);
    var intShape02 = venn.append("g")
        .append("path")
        .attr("d", function() {
          return "M" + intPt02[0] + "," + intPt02[2] + "A" + rad + "," + rad +
            " 0 0,1 " + intPt02[1] + "," + intPt02[3]+ "A" + rad + "," + rad +
            " 0 0,1 " + intPt02[0] + "," + intPt02[2];
        })
        .style("opacity", 1)
        .attr("id", vennId + "intShape02")
        .style('fill', color(4))
        .on("click", function(){
          $("table").remove();
          makeTable(int02NoMid, vennId + "int02")
          d3.select("#" + vennId + "int02").transition().duration(400).style("opacity", 1);
        })
        .on({
          "mouseover": function() {
            d3.select("#" + vennId + "intShape02").transition()
            .duration(300).style("opacity", 0.5); // slice
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function() {
            d3.select("#" + vennId + "intShape02").transition()
            .duration(300).style("opacity", 1); // slice opacity
            d3.select(this).style("cursor", "default");
          }
        });


    var botLeftCount = venn.append("text")
          .attr("transform", "translate(40,310)")
          .attr("text-anchor", "middle")
          .attr("fill", "white")
          .style("font-size", "20px")
          .text(int02.length-middle.length);

    //intersection12
    var int12 = findIntersection(names[1], names[2]);
    var int12NoMid = int12.diff(middle);
    var intPt12 = intersection(rightX, topY, rad, botMid, botY, rad);
    var intShape12 = venn.append("g")
        .append("path")
        .attr("id", vennId + "intShape12")
        .attr("d", function() {
          return "M" + intPt12[0] + "," + intPt12[2] + "A" + rad + "," + rad +
            " 0 0,1 " + intPt12[1] + "," + intPt12[3]+ "A" + rad + "," + rad +
            " 0 0,1 " + intPt12[0] + "," + intPt12[2];
        })
        .style("opacity", 1)
        .style('fill', color(5))
        .on("click", function(){
          $("table").remove();
          makeTable(int12NoMid, vennId + "int12")
          d3.select("#" + vennId + "int12").transition().duration(400).style("opacity", 1);
        })
        .on({
          "mouseover": function() {
            d3.select("#" + vennId + "intShape12").transition()
            .duration(300).style("opacity", 0.5); // slice
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function() {
            d3.select("#" + vennId + "intShape12").transition()
            .duration(300).style("opacity", 1); // slice opacity
            d3.select(this).style("cursor", "default");
          }
        });

    var botRightCount = venn.append("text")
          .attr("transform", "translate(270,310)")
          .attr("text-anchor", "middle")
          .attr("fill", "white")
          .style("font-size", "20px")
          .text(int12.length-middle.length);

    //intersection middle
    var intShape012 = venn.append("g")
      .append("path")
      .attr("d", function() {
        return "M" + intPt02[1] + "," + intPt02[3] + "A" + rad + "," + rad +
          " 0 0,1 " + intPt01[0] + "," + intPt01[2] + "A" + rad + "," + rad +
          " 0 0,1 " + intPt12[0] + "," + intPt12[2] + "A" + rad + "," + rad +
          " 0 0,1 " + intPt02[1] + "," + intPt02[3];
      })
      .style("opacity", 1)
      .attr("id", vennId + "intShape012")
      .style('fill', color(6))
      .on("click", function(){
        $("table").remove();
        makeTable(middle, vennId + "int012")
        d3.select("#" + vennId + "int012").transition().duration(400).style("opacity", 1);
      })
      .on({
        "mouseover": function() {
          d3.select("#" + vennId + "intShape012").transition()
          .duration(300).style("opacity", 0.5); // slice
          d3.select(this).style("cursor", "pointer");
        },
        "mouseout": function() {
          d3.select("#" + vennId + "intShape012").transition()
          .duration(300).style("opacity", 1); // slice opacity
          d3.select(this).style("cursor", "default");
        }
      });

    var midCount = venn.append("text")
          .attr("transform", "translate(155,260)")
          .attr("text-anchor", "middle")
          .attr("fill", "white")
          .style("font-size", "20px")
          .text(middle.length);

    var circle2Count = venn.append("text")
        .attr("transform", "translate(155,440)")
        .attr("text-anchor", "middle")
        .attr("fill", "white")
        .style("font-size", "20px")
        .text(names[2].length - int02.length - (int12.length-middle.length))

    //changing previous counts from when it wasn't 3+ circles
    topMidCount.attr("transform", "translate(155,140)")
      .text(int01.length-middle.length)
    circle0Count.attr("transform", "translate(-30,140)")
      .text(names[0].length - int01.length - (int02.length-middle.length))
    circle1Count.attr("transform", "translate(330,140)")
      .text(names[1].length - int01.length - (int12.length-middle.length))

    threeCNoMid0 = (names[0].diff(int01)).diff(int02NoMid);
    threeCNoMid1 = (names[1].diff(int12)).diff(int01NoMid);
    threeCNoMid2 = (names[2].diff(int12)).diff(int02NoMid);

    circle2.on("click", function(){
      $("table").remove();
      makeTable(threeCNoMid2, vennId + "cTable2")
      d3.select("#" + vennId + "cTable2").transition().duration(400).style("opacity", 1);
    })

  }

  //adding click for circle 0 and 1
  var twoCNoMid0 = names[0].diff(int01);
  var twoCNoMid1 = names[1].diff(int01);


  circle0.on("click", function(){
    $("table").remove();
    if (numCircles == 2) {
      makeTable(twoCNoMid0, vennId + "cTable0")
    } else if (numCircles == 3){
      makeTable(threeCNoMid0, vennId + "cTable0")
    }
    d3.select("#" + vennId + "cTable0").transition().duration(400).style("opacity", 1);
  })

  circle1.on("click", function(){
    $("table").remove();
    if (numCircles == 2) {
      makeTable(twoCNoMid1, vennId + "cTable1")
    } else if (numCircles == 3){
      makeTable(threeCNoMid1, vennId + "cTable1")
    }
    d3.select("#" + vennId + "cTable1").transition().duration(400).style("opacity", 1);
  })
}

// make for every clickable part of the venn diagram
function makeTable(names, tableId) {
  // makes however many rows, and 3 columns
  // also makes sure there is enough data to populate each subarray
  var tableData = [];
  var addElems = names.length % 3;
  names_copy = names.slice(0);
  for (var i = 0; i < 3-addElems; i++) {
    names_copy.push(" ");
  }

  //make multidimensional array of 10
  for (var i = 0; i < names.length; i += 3) {
    var temp = []
    for (var j = 0; j < 3; j++) {
      temp.push({"name": names_copy[i+j]})
    }
    tableData.push(temp)
  }


  var table = d3.select("#table_container").append('table')
    .attr("id", tableId)
    .style("border", "1px solid silver")
    .style("opacity", "0");

  var tbody = table.append("tbody");

  //create table rows
  var tr = tbody.selectAll("tr")
    .data(tableData)
    .enter()
    .append("tr")
    .style("border", "1px solid silver");

  // create table cells
  var td = tr.selectAll("td")
    .data(function(d) {return d;})
    .enter()
    .append("td")
    .style("font-size", "15px")
    .style("border", "1px solid silver")
    .attr("width", 250)
    .style("padding-top", "5px")
    .html(function(d,i){
        return "<a href=\"" + "/search?q=" + d.name + "\">" + d.name.replaceAll("_", " ") + "</a>"
      });

}


// makes an upset for 4+ sets, only goes up to 5 so far
function makeUpset(svg, numSets, names, vennId, datasets) {

  // make a group for the upset circle intersection things
  var upsetCircles = svg.append("g")
    .attr("id", "upsetCircles" + vennId)
    .attr("transform", function() {
      if (numSets > 4) {
        return "translate(-90, 400)"
      } else {
        return "translate(-150, 400)"
      }
    })
    .style("opacity", 0)
    .attr("visibility", "hidden");

  var rad = 13,
  height = 400;

  var width;
  if (numSets > 4) {
    width = 1100
  } else {
    width = 670
  }

  // computes intersections UNMANUALLY! goes up to  5 sets
  var data2 = []
  for (var i = 0; i < numSets; i++) {
    for (var j = i + 1; j < numSets; j++) {
      var temp = {}
      temp["set"] = i.toString() + j.toString()
      temp["names"] = findIntersection(names[i], names[j])
      data2.push(temp)
      for (var k = j + 1; k < numSets; k++) {
        var temp = {}
        temp["set"] = i.toString() + j.toString() + k.toString()
        temp["names"] = findIntersection(findIntersection(names[i], names[j]), names[k])
        data2.push(temp)
        for (var l = k + 1; l < numSets; l++) {
          var temp = {}
          temp["set"] = i.toString() + j.toString() + k.toString() + l.toString()
          temp["names"] = findIntersection(findIntersection(names[i], names[j]), findIntersection(names[k], names[l]))
          data2.push(temp)
          for (var m = l + 1; m < numSets; m++) {
            var temp = {}
            temp["set"] = i.toString() + j.toString() + k.toString() + l.toString() + m.toString()
            temp["names"] = findIntersection((findIntersection(names[i], names[j]), findIntersection(names[k], names[l])), names[m])
            data2.push(temp)
          }
        }
      }
    }
  }

  // makes sure data is unique
  var unique = []
  var newData = []
  for (var i = 0; i < data2.length; i++) {
    if (unique.indexOf(data2[i].set) == -1) {
      unique.push(data2[i].set)
      newData.push(data2[i])
    }
  }

  var data = newData

  // makes JSON objects for easy data parsing
  // also makes dataset labels
  for (var i = 0; i < numSets; i++) {
    data.push({"set": i.toString(), "names": names[i]})

    upsetCircles.append("text")
      .attr("dx", -20)
      .attr("dy", 5 + i * (rad*2.7))
      .attr("text-anchor", "end")
      .attr("fill", "black")
      .style("font-size", 13)
      .text(datasets[i])
  }

  // sort data decreasing
  data.sort(function(a, b) {
    return parseFloat(b.names.length) - parseFloat(a.names.length);
});

  // make the bars
  var upsetBars = svg.append("g")
    .attr("id", "upsetBars" + vennId)
    .attr("transform", function () {
      if (numSets > 4) {
        return "translate(-110, -30)"
      } else {
        return "translate(-170, -30)"
      }
    })
    .style("opacity", 0)
    .attr("visibility", "hidden");

    var nums = []
    for (var i = 0; i < data.length; i++) {
      nums.push(data[i].names.length)
    }

    var names = []
    for (var i = 0; i < data.length; i++) {
      names.push(data[i].names)
    }

  //set range for data by domain, and scale by range
  var xrange = d3.scale.linear()
    .domain([0, nums.length])
    .range([0, width]);


  var yrange = d3.scale.linear()
    .domain([0, nums[0]])
    .range([height, 0]);


  //set axes for graph
  var xAxis = d3.svg.axis()
    .scale(xrange)
    .orient("bottom")
    .tickPadding(2)
    .tickFormat(function(d,i) { return data[i].set})
    .tickValues(d3.range(data.length));

  var yAxis = d3.svg.axis()
    .scale(yrange)
    .orient("left")
    .tickSize(5)

  //add X axis
  upsetBars.append("g")
      .attr("class", "x axis")
      .attr("transform", "translate(0," +  height + ")")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .call(xAxis)
      .selectAll(".tick")
      .remove()


  // Add the Y Axis
  upsetBars.append("g")
      .attr("class", "y axis")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .call(yAxis)
      .selectAll("text")
      .attr("fill", "black")
      .attr("stroke", "none");

      // Y axis label
      upsetBars.append("text")
          .attr("text-anchor", "middle")
          .attr("fill","black")
          .attr("transform", "translate("+ -50 +","+(height/2)+")rotate(-90)")
          .text("# " + vennId.replace("_", " "));


  var chart = upsetBars.append('g')
  		.attr("transform", "translate(1,0)")
  		.attr('id','chart' + vennId)
      .style("opacity", 0)
      .attr("visibility", "hidden");

  // adding each bar
  chart.selectAll('.bar' + vennId)
  		.data(data)
  		.enter()
      .append('rect')
        .attr("class", "bar" + vennId)
        .attr('width', 15)
        .attr({
          'x':function(d,i){ return (rad-1) + i * (rad*2.7)},
          'y':function(d){ return yrange(d.names.length)}
        })
    		.style('fill', "darkslategrey")
    		.attr('height',function(d){ return height - yrange(d.names.length); })
        .on({
          "mouseover": function(d,i) {
            d3.select(this).transition().duration(300).style("opacity", 0.6);
            d3.select("#set" + i + "num" + vennId).transition().duration(300).style("opacity", 1);
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function(d,i) {
            d3.select(this).transition().duration(300).style("opacity", 1);
            d3.select("#set" + i + "num" + vennId).transition().duration(300).style("opacity", 0);
            d3.select(this).style("cursor", "default");
          }
        })
        .on("click", function(d,i){
          $("table").remove();
          console.log(d.length)
          makeTable(d.names, vennId + "uTable" + d.set)
          d3.select("#" + vennId + "uTable" + d.set).transition().duration(400).style("opacity", 1);
        })

  // for some reason i can't append text onto the end of the bars
  // and so i must do it here
  for (var i = 0; i < nums.length; i++) {
    upsetBars.append("text")
      .attr({'x': i * (rad*3) + 20, 'y': yrange(nums[i])-5})
      .attr("id", "set" + i + "num" + vennId)
      .style("text-anchor", "middle")
      .style("font-size", "13px")
      .attr("fill", "black")
      .style("opacity", 0)
      .text(nums[i])
    }

  //circles
  for (var i = 0; i < data.length; i++) {
    for (var j = 0; j < numSets; j++) {
      upsetCircles.append("circle")
        .attr("cx", i * (rad*2.7))
        .attr("cy", j * (rad*2.7))
        .attr("r", rad)
        .attr("id", vennId + "set" + i)
        .style("opacity", 1)
        .attr("fill", function() {
          if (data[i].set.indexOf(j.toString()) != -1) {
            return "darkslategrey"
          } else {
            return "silver"
          }
        })

    }

    if (data[i].set.length != 1) {
      upsetCircles.append("line")
        .attr("x1", i * (rad*2.7))
        .attr("y1", data[i].set.substring(0, 1) * (rad*2.7))
        .attr("x2", i * (rad*2.7))
        .attr("y2", data[i].set.substring(data[i].set.length - 1, data[i].set.length) * (rad*2.7))
        .style("stroke", "darkslategrey")
        .attr("stroke-width", 4)
    }
  }



}


function makeVenn(cell_lines, tissues, drugs, datasets, middle, plotId) { // names: [[],[]]
  //number of circles to make
  var numCircles = datasets.length;


  //position and dimensions
  var margin = {
    top: 200,
    right: 200,
    bottom: 100,
    left: 260
  };
  var width = 600;
  var height;
  if (numCircles == 2) {
    height = 300;
  } else if (numCircles == 3){
    height = 500;
  } else {
    height = 300;
  }

  //var color = d3.scale.category10();
  // var arrColor = ["#351431", "#351431", "#217675", "#351431", "#217675", "#217675", "#f1a629"]
  var purple = "#D35A62"//"#351431"
  var yellow = "#81ACD5"//"#d89017"
  var orange = "#79CD9E"//"#c64a2b"

  var arrColor = [purple, purple, yellow, purple, yellow, yellow, orange]

  var color = function(i) {
    return arrColor[i]
  }

  updateDimensions(window.innerWidth, margin);



  // make the canvas
  var svg = d3.select("#" + plotId) //"#" + plotId
      .append("svg")
          .attr("width", width + margin.left + margin.right)
          .attr("height", height + margin.top + margin.bottom)
      .attr("id", "venn" + plotId)
      .append("g")
          .attr("transform",
                "translate(" + margin.left + "," + margin.top + ")")
      .attr("fill", "white");


  // for resizing
  d3.select("#venn" + plotId)
  .attr( 'preserveAspectRatio',"xMinYMin meet")
  .attr("viewBox", function() {
    if (numCircles > 4) {
      return "0 0 1300 400"
    } else {
      return "0 0 850 400"
    }
  })
  .attr('width', function() {
    if (numCircles > 4) {
      return '850'
    } else {
      return '700'
    }
  })
  .attr("transform", function() {
    if (numCircles > 4) {
      return "translate(-70, 0)"
    } else {
      return "translate(0,0)"
    }
  })



  // graph title
  var graphTitle = svg.append("text")
    .attr("text-anchor", "middle")
    .attr("fill","black")
    .style("font-size", "30px")
    .attr("transform", function (){
      if (numCircles > 4) {
        return "translate("+ (width/2 + 200) +","+ -100 +")"
      } else {
        return "translate("+ (width/2 - 200) +","+ -100 +")"
      }
    });

  //"Show" title
  svg.append("text")
    .attr("x", function (){
      if (numCircles > 4) {
        return width + 290
      } else {
        return width-150
      }})
    .attr("y", -10)
    .attr("text-anchor", "start")
    .style("font-size", 20)
    .style("font-weight", "bold")
    .style("opacity", 1)
    .attr("fill", "black")
    .text("Show:")

  //determining the view of different venns
  var vennLabels = ["Cell Lines", "Tissues", "Compounds"]
  var vennIds = ["cell_lines", "tissues", "drugs"]

  //make all venns first
  if (numCircles <= 3) {
    makeCircles(svg, color, arrColor, numCircles, cell_lines, vennIds[0], datasets, middle[0]);
    makeCircles(svg, color, arrColor, numCircles, tissues, vennIds[1], datasets, middle[1]);
    makeCircles(svg, color, arrColor, numCircles, drugs, vennIds[2], datasets, middle[2]);
  } else {
    makeUpset(svg, numCircles, cell_lines, vennIds[0], datasets);
    makeUpset(svg, numCircles, tissues, vennIds[1], datasets);
    makeUpset(svg, numCircles, drugs, vennIds[2], datasets);
  }



  //nesting so that each element of Show view has an active property
  var nest = d3.nest()
    .key(function(d) {
        return d;
    })
    .entries([0,1,2]);

  nest.forEach(function(d,i) {
    var selection = svg.append("text")
      .attr("x", function (){
        if (numCircles > 4) {
          return width + 300
        } else {
          return width-120
        }})
      .attr("y", 25 + i * 30)
      .attr("id", "label" + vennIds[i])
      .attr("text-anchor", "start")
      .style("font-size", 18)
      .attr("fill", "silver");


    //default
    var gtitle = vennLabels[0] + ": ";
    for (var j = 0; j < datasets.length-1; j++) {
      gtitle = gtitle + datasets[j] + " + "
    }
    gtitle = gtitle + datasets[datasets.length-1]
    graphTitle.text(gtitle);

    d3.select("#labelcell_lines").attr("fill", "black");
    if (numCircles <= 3) {
      d3.select("#venn" + vennIds[0]).style("opacity", 0.8);
      d3.select("#venn" + vennIds[0]).attr("visibility", "visible");
    } else {
      d3.select("#upsetBars" + vennIds[0]).style("opacity", 1);
      d3.select("#upsetBars" + vennIds[0]).attr("visibility", "visible");
      d3.select("#upsetCircles" + vennIds[0]).style("opacity", 1);
      d3.select("#upsetCircles" + vennIds[0]).attr("visibility", "visible");
      d3.select("#chart" + vennIds[0]).style("opacity", 1);
      d3.select("#chart" + vennIds[0]).attr("visibility", "visible");
    }


    //on click of Show views
    selection.on("click", function(){
        var active   = active ? false : true;
        if (active) {
          var title = vennLabels[d.key] + ": ";
          for (var j = 0; j < datasets.length-1; j++) {
            title = title + datasets[j] + " + "
          }
          title = title + datasets[datasets.length-1]
          graphTitle.text(title);

          d3.select("#" + "label" + vennIds[d.key]).transition().duration(500).attr("fill", "black");
          if (numCircles <= 3) {
            d3.select("#venn" + vennIds[d.key]).attr("visibility", "visible");
            d3.select("#venn" + vennIds[d.key]).transition().duration(500).style("opacity", 0.8);
          } else {
            d3.select("#upsetBars" + vennIds[d.key]).attr("visibility", "visible");
            d3.select("#upsetBars" + vennIds[d.key]).transition().duration(500).style("opacity", 1);
            d3.select("#upsetCircles" + vennIds[d.key]).attr("visibility", "visible");
            d3.select("#upsetCircles" + vennIds[d.key]).transition().duration(500).style("opacity", 1);
            d3.select("#chart" + vennIds[d.key]).attr("visibility", "visible");
            d3.select("#chart" + vennIds[d.key]).transition().duration(500).style("opacity", 1);

          }

          for (var j = 0; j < 3; j++) {
            if (vennIds[d.key] != vennIds[j]) {
              d3.select("#" + "label" + vennIds[j]).transition().duration(500).attr("fill", "silver");
              if (numCircles <= 3) {
                d3.select("#venn" + vennIds[j]).transition().duration(500).style("opacity", 0);
                d3.select("#venn" + vennIds[j]).attr("visibility", "hidden");

              } else {
                d3.select("#upsetBars" + vennIds[j]).transition().duration(500).style("opacity", 0);
                d3.select("#upsetBars" + vennIds[j]).attr("visibility", "hidden");
                d3.select("#upsetCircles" + vennIds[j]).transition().duration(500).style("opacity", 0);
                d3.select("#upsetCircles" + vennIds[j]).attr("visibility", "hidden");
                d3.select("#chart" + vennIds[j]).transition().duration(500).style("opacity", 0);
                d3.select("#chart" + vennIds[j]).attr("visibility", "hidden");


              }
            }
          }
        } else {
          graphTitle.text();
          d3.select("#" + "label" + vennIds[d.key]).transition().duration(500).attr("fill", "silver");
          if (numCircles <= 3) {
            d3.select("#venn" + vennIds[d.key]).transition().duration(500).style("opacity", 0);
            d3.select("#venn" + vennIds[d.key]).attr("visibility", "hidden");

          } else {
            d3.select("#upsetBars" + vennIds[j]).transition().duration(500).style("opacity", 0);
            d3.select("#upsetBars" + vennIds[j]).attr("visibility", "hidden");
            d3.select("#upsetCircles" + vennIds[j]).transition().duration(500).style("opacity", 0);
            d3.select("#upsetCircles" + vennIds[j]).attr("visibility", "hidden");
            d3.select("#chart" + vennIds[j]).transition().duration(500).style("opacity", 0);
            d3.select("#chart" + vennIds[j]).attr("visibility", "hidden");

          }
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
      .text(vennLabels[i]);
    });

    var descriptor = "";
    for (var i = 0; i < datasets.length-1; i++) {
      descriptor = descriptor + datasets[i] + "-"
    }
    descriptor = descriptor + datasets[datasets.length-1]

    d3.select("#" + plotId).append("button")
        .attr("type","button")
        .attr("id", "button" + plotId)
        .attr("class", "downloadButton")
        .text("Download SVG")
        .on("click", function() {
            // download the svg
            downloadSVG(plotId, descriptor);
        });
}
