
function showError(error) {
  const errorNode = document.querySelector('#error');
  if (errorNode.firstChild) {
    errorNode.removeChild(errorNode.firstChild);
  }
  errorNode.appendChild(document.createTextNode(error.message || error));
}

function loadScreens() {
  return fetch('/api/screens', {
    method: 'GET',
    headers: {
      'Accepts': 'application/json'
    }
  }).then(res => {
    return res.json();
  }).then(payload => {
    return payload.content;
  }).catch(showError);
}

function startSession(offer, screen) {
  return fetch('/api/session', {
    method: 'POST',
    body: JSON.stringify({
      offer,
      screen
    }),
    headers: {
      'Content-Type': 'application/json'
    }
  }).then(res => {
    return res.json();
  }).then(payload => {
    return payload.content;
  }).then(msg => {
    return msg.answer;
  });
}

function createOffer(pc, { audio, video }) {
  return new Promise((accept, reject) => {
    pc.onicecandidate = evt => {
      if (!evt.candidate) {
        
        // ICE Gathering finished 
        const { sdp: offer } = pc.localDescription;
        accept(offer);
      }
    };
    pc.createOffer({
      offerToReceiveAudio: audio,
      offerToReceiveVideo: video
    }).then(ld => {
      pc.setLocalDescription(ld)
    }).catch(reject)
  });
}

function startRemoteSession(screen, remoteVideoNode, stream) {
  let pc;
  let dc;
  return Promise.resolve().then(() => {
    pc = new RTCPeerConnection({
      //iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
      //iceServers: [{ urls: 'stun:stun.ideasip.com' }]
      iceServers: [{ 
                     urls: 'stun:121.89.193.35:23478'
                   }, 
                   {
                     urls: 'turn:121.89.193.35:23478', 
                     username: 'wadahana', 
                     credential: '123456', 
                     credentialType: 'password',
                   }]
    });
    pc.ontrack = (evt) => {
      console.info('ontrack triggered');
      
      remoteVideoNode.srcObject = evt.streams[0];
      remoteVideoNode.play();
    };

    stream && stream.getTracks().forEach(track => {
      console.log('track')
      pc.addTrack(track, stream);
    })
    dc = pc.createDataChannel('event-channel');
    dc.binaryType = "arraybuffer";
    dc.onopen = (ev) => {
      console.log('dc open');
    }
    dc.onmessage = (ev) => {
      console.log('dc recv:' +ev.data);
    }
    dc.onclose = (ev) => {
      console.log('dc close');
    }
    return createOffer(pc, { audio: false, video: true });
  }).then(offer => {
    console.info(offer);
    return startSession(offer, screen);
  }).then(answer => {
    console.info(answer);
    return pc.setRemoteDescription(new RTCSessionDescription({
      sdp: answer,
      type: 'answer'
    }));
  }).then(() => {return {peer: pc, channel: dc};});
}

let peerConnection = null;
let dataChannel = null;
document.addEventListener('DOMContentLoaded', () => {
  
  let selectedScreen = 0;
  const remoteVideo = document.querySelector('#remote-video');
  const screenSelect = document.querySelector('#screen-select');
  const startStop = document.querySelector('#start-stop');
  
  loadScreens().then(response => {
    while (screenSelect.firstChild) {
      screenSelect.removeChild(screenSelect.firstChild);
    }
    screenSelect.appendChild(document.createElement('option'));
    response.screens.forEach(screen => {
      const option = document.createElement('option');
      //option.appendChild(document.createTextNode('Screen ' + (screen.index + 1)));
      option.appendChild(document.createTextNode(screen.name));
      option.setAttribute('value', screen.name);
      screenSelect.appendChild(option);
    });
  }).catch(showError);

  screenSelect.addEventListener('change', evt => {
    selectedScreen = evt.currentTarget.value//parseInt(evt.currentTarget.value, 10);
  });

  const enableStartStop = (enabled) => {
    if (enabled) {
      startStop.removeAttribute('disabled');
    } else {
      startStop.setAttribute('disabled', '');
    }
  }

  const setStartStopTitle = (title) => {
    startStop.removeChild(startStop.firstChild);
    startStop.appendChild(document.createTextNode(title));
  }

  startStop.addEventListener('click', () => {
    enableStartStop(false);

    //const userMediaPromise =  (adapter.browserDetails.browser === 'safari') ?
    //  navigator.mediaDevices.getUserMedia({ video: true }) : 
    //  Promise.resolve(null);
    const userMediaPromise = Promise.resolve(null);
    if (!peerConnection) {
      userMediaPromise.then(stream => {
        return startRemoteSession(selectedScreen, remoteVideo, stream).then(v => {
          remoteVideo.style.setProperty('visibility', 'visible');
          peerConnection = v.peer;
          dataChannel = v.channel;
        }).catch(showError).then(() => {
          enableStartStop(true);
          setStartStopTitle('Stop');
        });
      })
    } else {
      peerConnection.close();
      peerConnection = null;
      enableStartStop(true);
      setStartStopTitle('Start');
      remoteVideo.style.setProperty('visibility', 'collapse');
    }
  });
});

window.addEventListener('beforeunload', () => {
  if (peerConnection) {
    peerConnection.close();
  }
})

function sendData(data) {
  if (dataChannel != null && dataChannel.readyState == 'open') {
    dataChannel.send(data);
  }
}

function startHookEvent() {
  const evBox = document.querySelector('#event-box');
  //const evBox = document.querySelector('#remote-video')

  function getClientPosition(box, ev) {
    let x = (ev.offsetX / box.clientWidth).toFixed(4);
    let y = (ev.offsetY / box.clientHeight).toFixed(4);

    return {
      x:  x > 1.0 ? 1.0000 : x,
      y:  y > 1.0 ? 1.0000 : y,
    };
  }

  function sendMouseEvent(type, btn, pos) {
    let buffer = new ArrayBuffer(12);
    let dv = new DataView(buffer, 0, 12);
    dv.setUint8(0, 1);
    dv.setUint8(1, type);
    dv.setUint16(2, btn);
    dv.setFloat32(4, pos.x);
    dv.setFloat32(8, pos.y);
    //let event = new Uint8Array(buffer);
    if (dataChannel != null && dataChannel.readyState == 'open') {
      dataChannel.send(buffer);
    }
  }

  function sendKeyboardEvent(press, keyCode) {
      let buffer = new ArrayBuffer(8);
      let dv = new DataView(buffer, 0, 8);
      dv.setUint8(0, 2);
      dv.setUint8(1, press);
      dv.setUint8(2, 0xFF);
      dv.setUint8(3, 0xFF);
      dv.setUint32(4, keyCode);
      if (dataChannel != null && dataChannel.readyState == 'open') {
        dataChannel.send(buffer);
      }
  }

  evBox.addEventListener('keydown', function(ev) {
    console.log('keydown: keyCode:' + ev.keyCode + ',charCode:' + ev.charCode + ',key:' + ev.key);
    sendKeyboardEvent(true, ev.keyCode)
  })
  evBox.addEventListener('keyup', function(ev) {
    console.log('keyup: keyCode:' + ev.keyCode + ',charCode:' + ev.charCode + ',key:' + ev.key);
    sendKeyboardEvent(false, ev.keyCode)
  })
  evBox.addEventListener('keypress', function(ev) {
    //console.log('keypress: keyCode:' + ev.keyCode + ',charCode:' + ev.charCode + ',key:' + ev.key);
  })    
  evBox.addEventListener('mousemove', function(ev) {
  //evBox.onmousemove = function(ev) {
      let pos = getClientPosition(evBox, ev);
      //console.log('x:' + pos.x + ',y:' + pos.y);
      sendMouseEvent(5, 0, pos)
  });

  evBox.addEventListener('mousedown', function(ev) {
      //if (ev.button == 0) {
          let pos = getClientPosition(evBox, ev);
          sendMouseEvent(2, ev.button, pos)
          //console.log('down x:' + pos.x + ',y:' + pos.y);
      //}
  });

  evBox.addEventListener('mouseup', function(ev) {
      //if (ev.button == 0) {
          let pos = getClientPosition(evBox, ev);
          sendMouseEvent(3, ev.button, pos)
          //console.log('up x:' + pos.x + ',y:' + pos.y);
      //}
  });  
  evBox.addEventListener('mousewheel', function(ev) {
      let pos = getClientPosition(evBox, ev);
      delta = ev.wheelDelta 
      if (delta > 0) {
        if (delta > 32767) {
          delta = 32767
        }
      } else if (delta < 0) {
        if (delta < -32768) {
          delta = -32768
        }
      }
      sendMouseEvent(4, delta, pos)
      //console.log('wheel x:' + pos.x + ',y:' + pos.y);
  })
}

startHookEvent();
