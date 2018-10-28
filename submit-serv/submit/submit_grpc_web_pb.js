/**
 * @fileoverview gRPC-Web generated client stub for submit
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.submit = require('./submit_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.submit.SubmitServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.submit.SubmitServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!proto.submit.SubmitServiceClient} The delegate callback based client
   */
  this.delegateClient_ = new proto.submit.SubmitServiceClient(
      hostname, credentials, options);

};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.submit.PostRequest,
 *   !proto.submit.PostResponse>}
 */
const methodInfo_Post = new grpc.web.AbstractClientBase.MethodInfo(
  proto.submit.PostResponse,
  /** @param {!proto.submit.PostRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.submit.PostResponse.deserializeBinary
);


/**
 * @param {!proto.submit.PostRequest} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.submit.PostResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.submit.PostResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.submit.SubmitServiceClient.prototype.post =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/submit.SubmitService/Post',
      request,
      metadata,
      methodInfo_Post,
      callback);
};


/**
 * @param {!proto.submit.PostRequest} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.submit.PostResponse>}
 *     The XHR Node Readable Stream
 */
proto.submit.SubmitServicePromiseClient.prototype.post =
    function(request, metadata) {
  return new Promise((resolve, reject) => {
    this.delegateClient_.post(
      request, metadata, (error, response) => {
        error ? reject(error) : resolve(response);
      });
  });
};


module.exports = proto.submit;

