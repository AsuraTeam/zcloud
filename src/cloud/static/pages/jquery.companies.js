/**
* Theme: Adminox Admin Template
* Author: Coderthemes
* Component: Companies
* 
*/
$( document ).ready(function() {
    
    var DrawSparkline = function() {
        $('#company-1').sparkline([0, 23, 43, 35, 44, 45, 56, 37, 40], {
            type: 'line',
            width: $('#company-1').width(),
            height: '80',
            chartRangeMax: 50,
            lineColor: '#32c861',
            fillColor: 'rgba(50, 200, 97, 0.1)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
            maxSpotColor: false,
            minSpotColor: false,
            spotColor: false,
            lineWidth: 2
        });

        $('#company-2').sparkline([0, 25, 48, 32, 36, 20, 85, 56, 36], {
            type: 'line',
            width: $('#company-2').width(),
            height: '80',
            chartRangeMax: 50,
            lineColor: '#32c861',
            fillColor: 'rgba(50, 200, 97, 0.1)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
            maxSpotColor: false,
            minSpotColor: false,
            spotColor: false,
            lineWidth: 2
        });

        $('#company-3').sparkline([0, 36, 85, 25, 24, 56, 24, 28, 32], {
            type: 'line',
            width: $('#company-3').width(),
            height: '80',
            chartRangeMax: 50,
            lineColor: '#32c861',
            fillColor: 'rgba(50, 200, 97, 0.1)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
            maxSpotColor: false,
            minSpotColor: false,
            spotColor: false,
            lineWidth: 2
        });

        $('#company-4').sparkline([21, 28, 30, 35, 44, 82, 30, 37, 40], {
            type: 'line',
            width: $('#company-4').width(),
            height: '80',
            chartRangeMax: 50,
            lineColor: '#32c861',
            fillColor: 'rgba(50, 200, 97, 0.1)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
            maxSpotColor: false,
            minSpotColor: false,
            spotColor: false,
            lineWidth: 2
        });

        $('#company-5').sparkline([32, 28, 35, 89, 10, 15, 25, 37, 45], {
            type: 'line',
            width: $('#company-5').width(),
            height: '80',
            chartRangeMax: 50,
            lineColor: '#32c861',
            fillColor: 'rgba(50, 200, 97, 0.1)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
            maxSpotColor: false,
            minSpotColor: false,
            spotColor: false,
            lineWidth: 2
        });

        $('#company-6').sparkline([10, 25, 35, 35, 65, 75, 56, 37, 40], {
            type: 'line',
            width: $('#company-6').width(),
            height: '80',
            chartRangeMax: 50,
            lineColor: '#32c861',
            fillColor: 'rgba(50, 200, 97, 0.1)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
            maxSpotColor: false,
            minSpotColor: false,
            spotColor: false,
            lineWidth: 2
        });

        $('#company-7').sparkline([0, 23, 43, 35, 44, 45, 56, 37, 40], {
            type: 'line',
            width: $('#company-7').width(),
            height: '80',
            chartRangeMax: 50,
            lineColor: '#32c861',
            fillColor: 'rgba(50, 200, 97, 0.1)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
            maxSpotColor: false,
            minSpotColor: false,
            spotColor: false,
            lineWidth: 2
        });

        $('#company-8').sparkline([8, 19, 31, 35, 44, 50, 32, 37, 40], {
            type: 'line',
            width: $('#company-8').width(),
            height: '80',
            chartRangeMax: 50,
            lineColor: '#32c861',
            fillColor: 'rgba(50, 200, 97, 0.1)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
            maxSpotColor: false,
            minSpotColor: false,
            spotColor: false,
            lineWidth: 2
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