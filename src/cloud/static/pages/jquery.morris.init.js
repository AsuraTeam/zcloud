
/**
* Theme: Adminox Admin Template
* Author: Coderthemes
* Morris Chart
*/

!function($) {
    "use strict";

    var MorrisCharts = function() {};

    //creates line chart
    MorrisCharts.prototype.createLineChart = function(element, data, xkey, ykeys, labels, opacity, Pfillcolor, Pstockcolor, lineColors) {
        Morris.Line({
          element: element,
          data: data,
          xkey: xkey,
          ykeys: ykeys,
          labels: labels,
          fillOpacity: opacity,
          pointFillColors: Pfillcolor,
          pointStrokeColors: Pstockcolor,
          behaveLikeLine: true,
          gridLineColor: '#eef0f2',
          hideHover: 'auto',
          lineWidth: '3px',
          pointSize: 0,
          preUnits: '$',
          resize: true, //defaulted to true
          lineColors: lineColors
        });
    },
    //creates area chart
    MorrisCharts.prototype.createAreaChart = function(element, pointSize, lineWidth, data, xkey, ykeys, labels, lineColors) {
        Morris.Area({
            element: element,
            pointSize: 0,
            lineWidth: 0,
            data: data,
            xkey: xkey,
            ykeys: ykeys,
            labels: labels,
            hideHover: 'auto',
            resize: true,
            gridLineColor: '#eef0f2',
            lineColors: lineColors
        });
    },
    //creates area chart with dotted
    MorrisCharts.prototype.createAreaChartDotted = function(element, pointSize, lineWidth, data, xkey, ykeys, labels, Pfillcolor, Pstockcolor, lineColors) {
        Morris.Area({
            element: element,
            pointSize: 3,
            lineWidth: 1,
            data: data,
            xkey: xkey,
            ykeys: ykeys,
            labels: labels,
            hideHover: 'auto',
            pointFillColors: Pfillcolor,
            pointStrokeColors: Pstockcolor,
            resize: true,
            smooth: false,
            gridLineColor: '#eef0f2',
            lineColors: lineColors
        });
    },
    //creates Bar chart
    MorrisCharts.prototype.createBarChart  = function(element, data, xkey, ykeys, labels, lineColors) {
        Morris.Bar({
            element: element,
            data: data,
            xkey: xkey,
            ykeys: ykeys,
            labels: labels,
            hideHover: 'auto',
            resize: true, //defaulted to true
            gridLineColor: '#eeeeee',
            barSizeRatio: 0.4,
            xLabelAngle: 35,
            barColors: lineColors
        });
    },
    //creates Stacked chart
    MorrisCharts.prototype.createStackedChart  = function(element, data, xkey, ykeys, labels, lineColors) {
        Morris.Bar({
            element: element,
            data: data,
            xkey: xkey,
            ykeys: ykeys,
            stacked: true,
            labels: labels,
            hideHover: 'auto',
            resize: true, //defaulted to true
            gridLineColor: '#eeeeee',
            barColors: lineColors
        });
    },
    //creates Donut chart
    MorrisCharts.prototype.createDonutChart = function(element, data, colors) {
        Morris.Donut({
            element: element,
            data: data,
            barSize: 0.4,
            resize: true, //defaulted to true
            colors: colors
        });
    },
    MorrisCharts.prototype.init = function() {

        //create line chart
        var $data  = [
             { y: '2008', a: 50, b: 0 },
            { y: '2009', a: 75, b: 50 },
            { y: '2010', a: 30, b: 80 },
            { y: '2011', a: 50, b: 50 },
            { y: '2012', a: 75, b: 10 },
            { y: '2013', a: 50, b: 40 },
            { y: '2014', a: 75, b: 50 },
            { y: '2015', a: 100, b: 70 }
          ];
        this.createLineChart('morris-line-example', $data, 'y', ['a', 'b'], ['Series A', 'Series B'],['0.1'],['#ffffff'],['#999999'], ['#5553ce ', '#f06292']);

        //creating area chart
        var $areaData = [
            { y: '2009', a: 10, b: 20 },
            { y: '2010', a: 75,  b: 65 },
            { y: '2011', a: 50,  b: 40 },
            { y: '2012', a: 75,  b: 65 },
            { y: '2013', a: 50,  b: 40 },
            { y: '2014', a: 75,  b: 65 },
            { y: '2015', a: 90, b: 60 }
        ];
        this.createAreaChart('morris-area-example', 0, 0, $areaData, 'y', ['a', 'b'], ['Series A', 'Series B'], ['#5553ce', "#bdbdbd"]);

        //creating area chart with dotted
        var $areaDotData = [
            { y: '2009', a: 10, b: 20 },
            { y: '2010', a: 75,  b: 65 },
            { y: '2011', a: 50,  b: 40 },
            { y: '2012', a: 75,  b: 65 },
            { y: '2013', a: 50,  b: 40 },
            { y: '2014', a: 75,  b: 65 },
            { y: '2015', a: 90, b: 60 }
        ];
        this.createAreaChartDotted('morris-area-with-dotted', 0, 0, $areaDotData, 'y', ['a', 'b'], ['Series A', 'Series B'],['#ffffff'],['#999999'], ['#297ef6', "#bdbdbd"]);

        //creating bar chart
        var $barData  = [
            { y: '2009', a: 100, b: 90 , c: 40 },
            { y: '2010', a: 75,  b: 65 , c: 20 },
            { y: '2011', a: 50,  b: 40 , c: 50 },
            { y: '2012', a: 75,  b: 65 , c: 95 },
            { y: '2013', a: 50,  b: 40 , c: 22 },
            { y: '2014', a: 75,  b: 65 , c: 56 },
            { y: '2015', a: 100, b: 90 , c: 60 }
        ];
        this.createBarChart('morris-bar-example', $barData, 'y', ['a', 'b', 'c'], ['Series A', 'Series B', 'Series C'], ['#c086f3','#65acff', '#7ed321']);

        //creating Stacked chart
        var $stckedData  = [
            { y: '2005', a: 45, b: 180, c: 100 },
            { y: '2006', a: 75,  b: 65, c: 80 },
            { y: '2007', a: 100, b: 90, c: 56 },
            { y: '2008', a: 75,  b: 65, c: 89 },
            { y: '2009', a: 100, b: 90, c: 120 },
            { y: '2010', a: 75,  b: 65, c: 110 },
            { y: '2011', a: 50,  b: 40, c: 85 },
            { y: '2012', a: 75,  b: 65, c: 52 },
            { y: '2013', a: 50,  b: 40, c: 77 },
            { y: '2014', a: 75,  b: 65, c: 90 },
            { y: '2015', a: 100, b: 90, c: 130 },
            { y: '2016', a: 80, b: 65, c: 95 }
        ];
        this.createStackedChart('morris-bar-stacked', $stckedData, 'y', ['a', 'b', 'c'], ['Series A', 'Series B','Series C'], ['#4489e4', '#78c350','#ebeff2']);

        //creating donut chart
        var $donutData = [
                {label: "Electricity", value: 12},
                {label: "Financial", value: 30},
                {label: "Markets", value: 20}
            ];
        this.createDonutChart('morris-donut-example', $donutData, ['#297ef6', '#353d4a','#f2f6f8']);
    },
    //init
    $.MorrisCharts = new MorrisCharts, $.MorrisCharts.Constructor = MorrisCharts
}(window.jQuery),

//initializing 
function($) {
    "use strict";
    $.MorrisCharts.init();
}(window.jQuery);