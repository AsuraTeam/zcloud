/**
* Theme: Adminox Dashboard
* Author: Coderthemes
* Profile
*/


$( document ).ready(function() {

  var DrawSparkline = function() {
    $('#dashboard-1').sparkline([40, 51, 43, 35, 44, 45, 49], {
      type: 'bar',
      height: '40',
      barWidth: '10',
      barSpacing: '4',
      barColor: '#6ad9c3'
    });

      $('#dashboard-2').sparkline([44, 45, 49, 40, 51, 43, 35], {
          type: 'bar',
          height: '40',
          barWidth: '10',
          barSpacing: '4',
          barColor: '#ff89bb'
      });

      $('#dashboard-3').sparkline([43, 35, 44, 45, 49, 40, 51], {
          type: 'bar',
          height: '40',
          barWidth: '10',
          barSpacing: '4',
          barColor: '#9aa1f2'
      });
  }

  DrawSparkline();

  var resizeChart;

  $(window).resize(function(e) {
    clearTimeout(resizeChart);
    resizeChart = setTimeout(function() {
      DrawSparkline();
    }, 300);
  });
});