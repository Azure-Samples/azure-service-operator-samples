package azvotes

var templ = `<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>{{.Title}}</title>
    <style>
    body {
        background-color:#F8F8F8;
    }
    
    div#container {
        margin-top:5%;
    }
    
    div#space {
        display:block;
        margin: 0 auto;
        width: 500px;
        height: 10px;
     
    }
    
    div#logo {
        display:block;
        margin: 0 auto;
        width: 500px;
        text-align: center;
        font-size:30px;
        font-family:Helvetica;  
        /*border-bottom: 1px solid black;*/   
    }
    
    div#form {
        padding: 20px;
        padding-right: 20px;
        padding-top: 20px;
        display:block;
        margin: 0 auto;
        width: 500px;
        text-align: center;
        font-size:30px;
        font-family:Helvetica;  
        border-bottom: 1px solid black;
        border-top: 1px solid black;
    }
    
    div#results {
        display:block;
        margin: 0 auto;
        width: 500px;
        text-align: center;
        font-size:30px;
        font-family:Helvetica;  
    }
    
    .button {
        background-color: #4CAF50; /* Green */
        border: none;
        color: white;
        padding: 16px 32px;
        text-align: center;
        text-decoration: none;
        display: inline-block;
        font-size: 16px;
        margin: 4px 2px;
        -webkit-transition-duration: 0.4s; /* Safari */
        transition-duration: 0.4s;
        cursor: pointer;
        width: 250px;
    }
    
    .button1 {
        background-color: white; 
        color: black; 
        border: 2px solid #008CBA;
    }
    
    .button1:hover {
        background-color: #008CBA;
        color: white;
    }
    .button2 {
        background-color: white;
        color: black;
        border: 2px solid #555555;
    }
    
    .button2:hover {
        background-color: #555555;
        color: white;
    }
    
    .button3 {
        background-color: white; 
        color: black; 
        border: 2px solid #f44336;
    }
    
    .button3:hover {
        background-color: #f44336;
        color: white;
    }
    </style>
    <script language="JavaScript">
        function send(form){
        }
    </script>

</head>
<body>
    <div id="container">
        <form id="form" name="form" action="/" method="post"><center>
        <div id="logo">{{.Title}}</div>
        <div id="space"></div>
        <div id="form">
        <button name="vote" value="{{.Button1}}" onclick="send()" class="button button1">{{.Button1}}</button>
        <button name="vote" value="{{.Button2}}" onclick="send()" class="button button2">{{.Button2}}</button>
        <button name="vote" value="reset" onclick="send()" class="button button3">Reset</button>
        <div id="space"></div>
        <div id="space"></div>
        <div id="results"> {{.Button1}} - {{ .Value1 }} | {{.Button2}} - {{ .Value2 }} </div> 
        </form>        
        </div>
    </div>     
</body>
</html>`
