const {SubmitRequest, SubmitResponse} = require('./submit_pb.js');
const {SubmitServiceClient} = require('./submit_grpc_web_pb.js');

var submitService = new SubmitServiceClient('http://localhost:8080');

var request = new SubmitRequest(
    {
        url: "url1",
        title: "title1",
        body: "body1"
});

/*submitService.Submit(request, {}, function(err, response) {
  
});*/