
Todo: 

- Refactor and tidy up. Started this with the tree but that all needs
  to be separate files etc.

- Name!
- Use 'log' everywhere instead of printing stuff etc.
- Be more specific about how long things took - speed is the USP!
- Pull everything out into config

- Benchmark (e.g. find changed tests in big branch x 100), then optimise from comments


Later: 

- mutation testing. shell out to babel to change a value, run all affected tests and expect a failure
  if no failure, that's a problem - add it to a report (show a nice report with a diff of the change etc.)