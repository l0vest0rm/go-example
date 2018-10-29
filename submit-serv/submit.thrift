
struct SubmitRequest {
  1: string uid,
  2: string key,
  3: string url,
  4: string title,
  5: string body,
}

struct SubmitResponse {
  1: i32 code = 0,
  2: string message,
}

service SubmitService {
  /**
   * A method definition looks like C code. It has a return type, arguments,
   * and optionally a list of exceptions that it may throw. Note that argument
   * lists and exception lists are specified using the exact same syntax as
   * field lists in struct or exception definitions.
   */
  void ping(),
  
  // submit
  SubmitResponse submit(1:SubmitRequest req),
}