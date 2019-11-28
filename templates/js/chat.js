 $(function(){

  var socket = null;
  var msgBox = $("#chatbox textarea");
  var messages = $("#messages");

  $("#chatbox").submit(function(){

    if (!msgBox.val()) return false;
    if (!socket) {
      alert("Błąd: Brak połączenia z serwerem.");
      return false;
    }

    socket.send(msgBox.val());
    msgBox.val("");
    return false;

  }); 

  if (!window["WebSocket"]) {
    alert("Błąd: Twoja przeglądarka nie obsługuje technologii WebScocket.")
  } else {
    socket = new WebSocket("ws://{{.Host}}/room");
    socket.onclose = function() {
      alert("Połączenie zostało zamknięte.");
    }
    socket.onmessage = function(e) {
      messages.append($("<li>").text(e.data));
    }
  }

});
