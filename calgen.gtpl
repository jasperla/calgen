<html>
  <head>
    <title>Kalender Generator</title>
    <link rel="stylesheet" href="style.css"></link>
  </head>

  <body>
    <form action='/calgen' method='post'>
      <fieldset>
	<legend>Parameters</legend>

	<table class="generator">
	  <tr>
	    <td><label>Begin datum</label></td>
	    <td><input type='date' name='begindate' required></td>
	  </tr>
	  <tr>
	    <td><label>Aantal weken</label></td>
	    <td><select name="weeks">
              <option value="1">1</option>
	      <option value="2">2</option>
	      <option value="3">3</option>
	      <option value="4">4</option>
  	      <option value="5">5</option>
  	      <option value="6">6</option>
              </select>
	    </td>
	  </tr>
	  </table>

	</p>

        <input type='submit' value='Genereer PDF'>
      </fieldset>
    </form>
  </body>
</html>
