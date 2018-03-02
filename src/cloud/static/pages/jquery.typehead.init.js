/**
 * Theme: Adminox Admin Template
 * Author: Coderthemes
 * Typeahead
 */



$(document).ready(function() {

  // the basics
  // ----------

  var substringMatcher = function(strs) {
    return function findMatches(q, cb) {
      var matches, substringRegex;

      // an array that will be populated with substring matches
      matches = [];

      // regex used to determine if a string contains the substring `q`
      substrRegex = new RegExp(q, 'i');

      // iterate through the pool of strings and for any string that
      // contains the substring `q`, add it to the `matches` array
      $.each(strs, function(i, str) {
        if (substrRegex.test(str)) {
          matches.push(str);
        }
      });

      cb(matches);
    };
  };

  var states = ['Alabama', 'Alaska', 'Arizona', 'Arkansas', 'California',
    'Colorado', 'Connecticut', 'Delaware', 'Florida', 'Georgia', 'Hawaii',
    'Idaho', 'Illinois', 'Indiana', 'Iowa', 'Kansas', 'Kentucky', 'Louisiana',
    'Maine', 'Maryland', 'Massachusetts', 'Michigan', 'Minnesota',
    'Mississippi', 'Missouri', 'Montana', 'Nebraska', 'Nevada', 'New Hampshire',
    'New Jersey', 'New Mexico', 'New York', 'North Carolina', 'North Dakota',
    'Ohio', 'Oklahoma', 'Oregon', 'Pennsylvania', 'Rhode Island',
    'South Carolina', 'South Dakota', 'Tennessee', 'Texas', 'Utah', 'Vermont',
    'Virginia', 'Washington', 'West Virginia', 'Wisconsin', 'Wyoming'
  ];

  $('#the-basics').typeahead({
    hint: true,
    highlight: true,
    minLength: 1
  },
  {
    name: 'states',
    source: substringMatcher(states)
  });

  // bloodhound
  // ----------

  // constructs the suggestion engine
  var states = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.whitespace,
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    // `states` is an array of state names defined in "The Basics"
    local: states
  });

  $('#bloodhound').typeahead({
    hint: true,
    highlight: true,
    minLength: 1
  },
  {
    name: 'states',
    source: states
  });

  // prefetch
  // --------

  var countries = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.whitespace,
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    // url points to a json file that contains an array of country names, see
    // https://github.com/twitter/typeahead.js/blob/gh-pages/data/countries.json
    prefetch: '../plugins/typeahead/data/countries.json'
  });

  // passing in `null` for the `options` arguments will result in the default
  // options being used
  $('#prefetch').typeahead(null, {
    name: 'countries',
    source: countries
  });

  // remote
  // ------

  var bestPictures = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('value'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    prefetch: '../plugins/typeahead/data/post_1960.json',
    remote: {
      url: '../plugins/typeahead/data/%QUERY.json',
      wildcard: '%QUERY'
    }
  });

  $('#remote').typeahead(null, {
    name: 'best-pictures',
    display: 'value',
    source: bestPictures
  });

  // default suggestions
  // -------------------

  var nflTeams = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('team'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    identify: function(obj) { return obj.team; },
    prefetch: '../plugins/typeahead/data/nfl.json'
  });

  function nflTeamsWithDefaults(q, sync) {
    if (q === '') {
      sync(nflTeams.get('Detroit Lions', 'Green Bay Packers', 'Chicago Bears'));
    }

    else {
      nflTeams.search(q, sync);
    }
  }

  $('#default-suggestions').typeahead({
    minLength: 0,
    highlight: true
  },
  {
    name: 'nfl-teams',
    display: 'team',
    source: nflTeamsWithDefaults
  });

  // custom templates
  // ----------------

  $('#custom-templates').typeahead(null, {
    name: 'best-pictures',
    display: 'value',
    source: bestPictures,
    templates: {
      empty: [
        '<div class="typeahead-empty-message">',
          'unable to find any Best Picture winners that match the current query',
        '</div>'
      ].join('\n'),
      suggestion: Handlebars.compile('<div><strong>{{value}}</strong> – {{year}}</div>')
    }
  });

  // multiple datasets
  // -----------------

  var nbaTeams = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('team'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    prefetch: '../plugins/typeahead/data/nba.json'
  });

  var nhlTeams = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('team'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    prefetch: '../plugins/typeahead/data/nhl.json'
  });

  $('#multiple-datasets').typeahead({
    highlight: true
  },
  {
    name: 'nba-teams',
    display: 'team',
    source: nbaTeams,
    templates: {
      header: '<h5 class="league-name">NBA Teams</h5>'
    }
  },
  {
    name: 'nhl-teams',
    display: 'team',
    source: nhlTeams,
    templates: {
      header: '<h5 class="league-name">NHL Teams</h5>'
    }
  });


});