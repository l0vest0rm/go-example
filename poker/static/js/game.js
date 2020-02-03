PG.Game = function (game) {

  this.roomId = 1;
  this.totalPlayers = 5;
  this.players = [];

  this.titleBar = null;
  this.tableId = 0;
  this.shotLayer = null;

  this.tablePoker = [];
  this.tablePokerPic = {};

  this.lastShotPlayer = null;

  this.whoseTurn = 0;
  this.uid = 1710;

};

PG.Game.prototype = {

  init: function (roomId) {
    this.roomId = roomId;
  },

  debug_log: function (obj) {
    console.log('*******');
    console.log(obj);
    console.log('********');
  },

  create: function () {
    this.stage.backgroundColor = '#182d3b';

    for (i = 0; i < this.totalPlayers; i++) {
      this.players.push(PG.createPlay(this, i));
    }

    this.players[0].updateInfo(PG.playerInfo.uid, PG.playerInfo.username);
    PG.Socket.connect(this.onopen.bind(this), this.onmessage.bind(this), this.onerror.bind(this));

    this.createTitleBar();
  },

  onopen: function () {
    console.log('socket onopen');
    PG.Socket.send({
      code: PG.Protocol.REQ_JOIN_ROOM,
      uid: this.uid,
      data: this.roomId
    });
  },

  onerror: function () {
    console.log('socket connect onerror');
  },

  send_message: function (request) {
    PG.Socket.send(request);
  },

  onmessage: function (msg) {
    var opcode = msg.code;
    switch (opcode) {
      case PG.Protocol.RSP_JOIN_ROOM:
        if (this.roomId == 1) {
          PG.Socket.send({
            code: PG.Protocol.REQ_JOIN_TABLE,
            tableId: -1,
            uid: this.uid,
          });
        } else {
          this.createTableLayer(msg.data);
        }
        break;
      case PG.Protocol.RSP_TABLE_LIST:
        this.createTableLayer(msg.data);
        break;
      case PG.Protocol.RSP_NEW_TABLE:
        this.tableId = msg.data;
        this.titleBar.text = '房间:' + this.tableId;
        break;
      case PG.Protocol.RSP_JOIN_TABLE:
        this.tableId = msg.tableId;
        this.titleBar.text = '房间:' + this.tableId;
        var playerIds = msg.data;

        for (var i = 0; i < playerIds.length; i++) {
          this.players[i].updateInfo(playerIds[i].uid, playerIds[i].name);
        }
        break;
      case PG.Protocol.YOUR_TURN:
        this.startPlay();
        break;
      case PG.Protocol.INVALID_POCKER:
        this.invalidPoker();
        break;
      case PG.Protocol.RSP_DEAL_POKER:
        var playerId = msg.uid;
        var pokers = msg.data;
        console.log(pokers);
        this.dealPoker(pokers);
        this.whoseTurn = this.uidToSeat(playerId);
        //this.startCallScore(0);
        //this.startPlay();
        break;
      case PG.Protocol.RSP_CALL_SCORE:
        var playerId = packet[1];
        var score = packet[2];
        var callend = packet[3];
        this.debug_log(callend);
        this.whoseTurn = this.uidToSeat(playerId);
        //this.debug_log(playerId);

        var hanzi = ['不叫', "一分", "两分", "三分"];
        this.players[this.whoseTurn].say(hanzi[score]);
        if (!callend) {
          this.whoseTurn = (this.whoseTurn + 1) % 5;
          this.startCallScore(score);
        }
        break;
      case PG.Protocol.RSP_SHOW_POKER:
        this.whoseTurn = this.uidToSeat(packet[1]);
        this.tablePoker[0] = packet[2][0];
        this.tablePoker[1] = packet[2][1];
        this.tablePoker[2] = packet[2][2];
        this.players[this.whoseTurn].setLandlord();
        this.showLastThreePoker();
        break;
      case PG.Protocol.RSP_SHOT_POKER:
        this.handleShotPoker(msg);
        break;
      case PG.Protocol.RSP_GAME_OVER:
        /*var winner = packet[1];
        var coin = packet[2];

        var loserASeat = this.uidToSeat(packet[3][0]);
        this.players[loserASeat].replacePoker(packet[3], 1);
        this.players[loserASeat].reDealPoker();

        var loserBSeat = this.uidToSeat(packet[4][0]);
        this.players[loserBSeat].replacePoker(packet[4], 1);
        this.players[loserBSeat].reDealPoker();
        //                 this.players[loserBSeat].removeAllPoker();
        //               this.players[loserASeat].pokerInHand = [];

        this.whoseTurn = this.uidToSeat(winner);*/

        function gameOver() {
          alert(msg.data);
          PG.Socket.send({
            code: PG.Protocol.REQ_RESTART,
            tableId: this.tableId
          });
          this.cleanWorld();
        }
        this.game.time.events.add(3000, gameOver, this);
        break;
      case PG.Protocol.RSP_CHEAT:
        var seat = this.uidToSeat(packet[1]);
        this.players[seat].replacePoker(packet[2], 0);
        this.players[seat].reDealPoker();
        break;
      case PG.Protocol.RSP_RESTART:
        this.restart();
      default:
        console.log("UNKNOWN PACKET:", packet)
    }
  },

  cleanWorld: function () {
    for (i = 0; i < this.players.length; i++) {
      this.players[i].cleanPokers();
      try {
        this.players[i].uiLeftPoker.kill();
      } catch (err) {}
      this.players[i].uiHead.frameName = 'icon_farmer.png';
    }

    for (var i = 0; i < this.tablePoker.length; i++) {
      var p = this.tablePokerPic[this.tablePoker[i]];
      // p.kill();
      p.destroy();
    }
  },

  restart: function () {
    this.players = [];
    this.shotLayer = null;

    this.tablePoker = [];
    this.tablePokerPic = {};

    this.lastShotPlayer = null;

    this.whoseTurn = 0;

    this.stage.backgroundColor = '#182d3b';

    for (i = 0; i < this.totalPlayers; i++) {
      this.players.push(PG.createPlay(this, i));
    }

    player_id = [1, 11, 12];
    for (var i = 0; i < 5; i++) {
      //this.players[i].uiHead.kill();
      this.players[i].updateInfo(player_id[i], ' ');
    }

    // this.send_message([PG.Protocol.REQ_DEAL_POKEER, -1]);
    //        PG.Socket.send([PG.Protocol.REQ_JOIN_TABLE, this.tableId]);
  },

  update: function () {},

  uidToSeat: function (uid) {
    for (var i = 0; i < this.players.length; i++) {
      //	        this.debug_log(this.players[i].uid);
      if (uid == this.players[i].uid)
        return i;
    }
    console.log('ERROR uidToSeat:' + uid);
    return -1;
  },

  dealPoker: function (pokers) {

    len = pokers.length
    for (var i = 0; i < len; i++) {
      this.players[4].pokerInHand.push(54);
      this.players[3].pokerInHand.push(54);
      this.players[2].pokerInHand.push(54);
      this.players[1].pokerInHand.push(54);
      this.players[0].pokerInHand.push(pokers.pop());
    }

    this.players[0].dealPoker();
    this.players[1].dealPoker();
    this.players[2].dealPoker();
    this.players[3].dealPoker();
    this.players[4].dealPoker();
    //this.game.time.events.add(1000, function() {
    //    this.send_message([PG.Protocol.REQ_CHEAT, this.players[1].uid]);
    //    this.send_message([PG.Protocol.REQ_CHEAT, this.players[2].uid]);
    //}, this);
  },

  showLastThreePoker: function () {
    for (var i = 0; i < 3; i++) {
      var pokerId = this.tablePoker[i];
      var p = this.tablePoker[i + 3];
      p.id = pokerId;
      p.frame = pokerId;
      this.game.add.tween(p).to({
        x: this.game.world.width / 2 + (i - 1) * 60
      }, 600, Phaser.Easing.Default, true);
    }
    this.game.time.events.add(1500, this.dealLastThreePoker, this);
  },

  dealLastThreePoker: function () {
    var turnPlayer = this.players[this.whoseTurn];

    for (var i = 0; i < 3; i++) {
      var pid = this.tablePoker[i];
      var poker = this.tablePoker[i + 3];
      turnPlayer.pokerInHand.push(pid);
      turnPlayer.pushAPoker(poker);
    }
    turnPlayer.sortPoker();
    if (this.whoseTurn == 0) {
      turnPlayer.arrangePoker();
      for (var i = 0; i < 3; i++) {
        var p = this.tablePoker[i + 3];
        var tween = this.game.add.tween(p).to({
          y: this.game.world.height - PG.PH * 0.8
        }, 400, Phaser.Easing.Default, true);

        function adjust(p) {
          this.game.add.tween(p).to({
            y: this.game.world.height - PG.PH / 2
          }, 400, Phaser.Easing.Default, true, 400);
        };
        tween.onComplete.add(adjust, this, p);
      }
    } else {
      var first = turnPlayer.findAPoker(54);
      for (var i = 0; i < 3; i++) {
        var p = this.tablePoker[i + 3];
        p.frame = 54;
        p.frame = 54;
        this.game.add.tween(p).to({
          x: first.x,
          y: first.y
        }, 200, Phaser.Easing.Default, true);
      }
    }

    this.tablePoker = [];
    this.lastShotPlayer = turnPlayer;
    if (this.whoseTurn == 0) {
      this.startPlay();
    }
  },

  handleShotPoker: function (msg) {
    this.whoseTurn = this.uidToSeat(msg.uid);
    var turnPlayer = this.players[this.whoseTurn];
    var pokers = msg.data;

    if (this.whoseTurn == 0) {
      this.players[0].hindButton();
    }

    if (pokers.length == 0) {
      this.players[this.whoseTurn].say("不出");
    } else {
      var pokersPic = {};
      pokers.sort(PG.Poker.comparePoker);
      var count = pokers.length;
      var gap = Math.min((this.game.world.width - PG.PW * 2) / count, PG.PW * 0.36);
      for (var i = 0; i < count; i++) {
        var p = turnPlayer.findAPoker(pokers[i]);
        p.id = pokers[i];
        p.frame = pokers[i];
        p.bringToTop();
        this.game.add.tween(p).to({
          x: this.game.world.width / 2 + (i - count / 2) * gap,
          y: this.game.world.height * 0.4
        }, 500, Phaser.Easing.Default, true);

        turnPlayer.removeAPoker(pokers[i]);
        pokersPic[p.id] = p;
      }

      for (var i = 0; i < this.tablePoker.length; i++) {
        var p = this.tablePokerPic[this.tablePoker[i]];
        // p.kill();
        p.destroy();
      }
      this.tablePoker = pokers;
      this.tablePokerPic = pokersPic;
      this.lastShotPlayer = turnPlayer;
      turnPlayer.arrangePoker();
    }

    this.whoseTurn = (this.whoseTurn + 1) % this.players.length;
    /*if (this.whoseTurn == 0 && this.players[this.whoseTurn].pokerInHand.length > 0) {
      this.game.time.events.add(1000, this.startPlay, this);
    }*/
  },

  startCallScore: function (minscore) {
    function btnTouch(btn) {
      this.send_message({
        code: PG.Protocol.REQ_CALL_SCORE,
        tableId: this.tableId,
        data: btn.score
      });
      btn.parent.destroy();
      var audio = this.game.add.audio('f_score_' + btn.score);
      audio.play();
    };

    if (this.whoseTurn == 0) {
      var step = this.game.world.width / 6;
      var ss = [1.5, 1, 0.5, 0];
      var sx = this.game.world.width / 2 - step * ss[minscore];
      var sy = this.game.world.height * 0.6;
      var group = this.game.add.group();
      var pass = this.game.make.button(sx, sy, "btn", btnTouch, this, 'score_0.png', 'score_0.png', 'score_0.png');
      pass.anchor.set(0.5, 0);
      pass.score = 0;
      group.add(pass);
      sx += step;

      for (var i = minscore + 1; i <= 3; i++) {
        var tn = 'score_' + i + '.png';
        var call = this.game.make.button(sx, sy, "btn", btnTouch, this, tn, tn, tn);
        call.anchor.set(0.5, 0);
        call.score = i;
        group.add(call);
        sx += step;
      }
    } else {
      // TODO show clock on player
    }

  },

  startPlay: function () {
    if (this.isLastShotPlayer()) {
      this.players[0].playPoker([]);
    } else {
      this.players[0].playPoker(this.tablePoker);
    }
  },

  invalidPoker: function () {
    this.players[0].invalidPoker();
  },

  finishPlay: function (pokers) {
    this.send_message({
      code: PG.Protocol.REQ_SHOT_POKER,
      tableId: this.tableId,
      uid: this.players[0].uid,
      data: pokers
    });
  },

  isLastShotPlayer: function () {
    return this.players[this.whoseTurn] == this.lastShotPlayer;
  },

  createTableLayer: function (tables) {
    tables.push([-1, 0]);

    var group = this.game.add.group();
    this.game.world.bringToTop(group);
    var gc = this.game.make.graphics(0, 0);
    gc.beginFill(0x00000080);
    gc.endFill();
    group.add(gc);
    var style = {
      font: "22px Arial",
      fill: "#fff",
      align: "center"
    };

    for (var i = 0; i < tables.length; i++) {
      var sx = this.game.world.width * (i % 6 + 1) / (6 + 1);
      var sy = this.game.world.height * (Math.floor(i / 6) + 1) / (4 + 1);

      var table = this.game.make.button(sx, sy, 'btn', this.onJoin, this, 'table.png', 'table.png', 'table.png');
      table.anchor.set(0.5, 1);
      table.tableId = tables[i][0];
      group.add(table);

      var text = this.game.make.text(sx, sy, '房间:' + tables[i][0] + '人数:' + tables[i][1], style);
      text.anchor.set(0.5, 0);
      group.add(text);

      if (i == tables.length - 1) {
        text.text = '新建房间';
      }
    }
  },

  quitGame: function () {
    this.state.start('MainMenu');
  },

  createTitleBar: function () {
    var style = {
      font: "22px Arial",
      fill: "#fff",
      align: "center"
    };
    this.titleBar = this.game.add.text(this.game.world.centerX, 0, '房间:', style);
  },

  onJoin: function (btn) {
    if (btn.tableId == -1) {
      this.send_message({
        code: PG.Protocol.REQ_NEW_TABLE
      });
    } else {
      this.send_message({
        code: PG.Protocol.REQ_JOIN_TABLE,
        tableId: btn.tableId
      });
    }
    btn.parent.destroy();
  }
};