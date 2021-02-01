// source: manifest.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require("google-protobuf");
var goog = jspb;
var global = Function("return this")();

goog.exportSymbol("proto.manifest.BuildInfo", null, global);
goog.exportSymbol("proto.manifest.HttpMethod", null, global);
goog.exportSymbol("proto.manifest.Manifest", null, global);
goog.exportSymbol("proto.manifest.Route", null, global);
goog.exportSymbol("proto.manifest.Version", null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.manifest.Manifest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.manifest.Manifest.repeatedFields_,
    null
  );
};
goog.inherits(proto.manifest.Manifest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.manifest.Manifest.displayName = "proto.manifest.Manifest";
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.manifest.BuildInfo = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.manifest.BuildInfo, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.manifest.BuildInfo.displayName = "proto.manifest.BuildInfo";
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.manifest.Version = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.manifest.Version, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.manifest.Version.displayName = "proto.manifest.Version";
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.manifest.Route = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.manifest.Route, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.manifest.Route.displayName = "proto.manifest.Route";
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.manifest.Manifest.repeatedFields_ = [4];

if (jspb.Message.GENERATE_TO_OBJECT) {
  /**
   * Creates an object representation of this proto.
   * Field names that are reserved in JavaScript and will be renamed to pb_name.
   * Optional fields that are not set will be set to undefined.
   * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
   * For the list of reserved names please see:
   *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
   * @param {boolean=} opt_includeInstance Deprecated. whether to include the
   *     JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @return {!Object}
   */
  proto.manifest.Manifest.prototype.toObject = function (opt_includeInstance) {
    return proto.manifest.Manifest.toObject(opt_includeInstance, this);
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.manifest.Manifest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.manifest.Manifest.toObject = function (includeInstance, msg) {
    var f,
      obj = {
        name: jspb.Message.getFieldWithDefault(msg, 1, ""),
        version:
          (f = msg.getVersion()) &&
          proto.manifest.Version.toObject(includeInstance, f),
        buildInfo:
          (f = msg.getBuildInfo()) &&
          proto.manifest.BuildInfo.toObject(includeInstance, f),
        routesList: jspb.Message.toObjectList(
          msg.getRoutesList(),
          proto.manifest.Route.toObject,
          includeInstance
        ),
      };

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.manifest.Manifest}
 */
proto.manifest.Manifest.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.manifest.Manifest();
  return proto.manifest.Manifest.deserializeBinaryFromReader(msg, reader);
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.manifest.Manifest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.manifest.Manifest}
 */
proto.manifest.Manifest.deserializeBinaryFromReader = function (msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
      case 1:
        var value = /** @type {string} */ (reader.readString());
        msg.setName(value);
        break;
      case 2:
        var value = new proto.manifest.Version();
        reader.readMessage(
          value,
          proto.manifest.Version.deserializeBinaryFromReader
        );
        msg.setVersion(value);
        break;
      case 3:
        var value = new proto.manifest.BuildInfo();
        reader.readMessage(
          value,
          proto.manifest.BuildInfo.deserializeBinaryFromReader
        );
        msg.setBuildInfo(value);
        break;
      case 4:
        var value = new proto.manifest.Route();
        reader.readMessage(
          value,
          proto.manifest.Route.deserializeBinaryFromReader
        );
        msg.addRoutes(value);
        break;
      default:
        reader.skipField();
        break;
    }
  }
  return msg;
};

/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.manifest.Manifest.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.manifest.Manifest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.manifest.Manifest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.manifest.Manifest.serializeBinaryToWriter = function (message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getVersion();
  if (f != null) {
    writer.writeMessage(2, f, proto.manifest.Version.serializeBinaryToWriter);
  }
  f = message.getBuildInfo();
  if (f != null) {
    writer.writeMessage(3, f, proto.manifest.BuildInfo.serializeBinaryToWriter);
  }
  f = message.getRoutesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.manifest.Route.serializeBinaryToWriter
    );
  }
};

/**
 * optional string name = 1;
 * @return {string}
 */
proto.manifest.Manifest.prototype.getName = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};

/**
 * @param {string} value
 * @return {!proto.manifest.Manifest} returns this
 */
proto.manifest.Manifest.prototype.setName = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional Version version = 2;
 * @return {?proto.manifest.Version}
 */
proto.manifest.Manifest.prototype.getVersion = function () {
  return /** @type{?proto.manifest.Version} */ (jspb.Message.getWrapperField(
    this,
    proto.manifest.Version,
    2
  ));
};

/**
 * @param {?proto.manifest.Version|undefined} value
 * @return {!proto.manifest.Manifest} returns this
 */
proto.manifest.Manifest.prototype.setVersion = function (value) {
  return jspb.Message.setWrapperField(this, 2, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.manifest.Manifest} returns this
 */
proto.manifest.Manifest.prototype.clearVersion = function () {
  return this.setVersion(undefined);
};

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.manifest.Manifest.prototype.hasVersion = function () {
  return jspb.Message.getField(this, 2) != null;
};

/**
 * optional BuildInfo build_info = 3;
 * @return {?proto.manifest.BuildInfo}
 */
proto.manifest.Manifest.prototype.getBuildInfo = function () {
  return /** @type{?proto.manifest.BuildInfo} */ (jspb.Message.getWrapperField(
    this,
    proto.manifest.BuildInfo,
    3
  ));
};

/**
 * @param {?proto.manifest.BuildInfo|undefined} value
 * @return {!proto.manifest.Manifest} returns this
 */
proto.manifest.Manifest.prototype.setBuildInfo = function (value) {
  return jspb.Message.setWrapperField(this, 3, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.manifest.Manifest} returns this
 */
proto.manifest.Manifest.prototype.clearBuildInfo = function () {
  return this.setBuildInfo(undefined);
};

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.manifest.Manifest.prototype.hasBuildInfo = function () {
  return jspb.Message.getField(this, 3) != null;
};

/**
 * repeated Route routes = 4;
 * @return {!Array<!proto.manifest.Route>}
 */
proto.manifest.Manifest.prototype.getRoutesList = function () {
  return /** @type{!Array<!proto.manifest.Route>} */ (jspb.Message.getRepeatedWrapperField(
    this,
    proto.manifest.Route,
    4
  ));
};

/**
 * @param {!Array<!proto.manifest.Route>} value
 * @return {!proto.manifest.Manifest} returns this
 */
proto.manifest.Manifest.prototype.setRoutesList = function (value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};

/**
 * @param {!proto.manifest.Route=} opt_value
 * @param {number=} opt_index
 * @return {!proto.manifest.Route}
 */
proto.manifest.Manifest.prototype.addRoutes = function (opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(
    this,
    4,
    opt_value,
    proto.manifest.Route,
    opt_index
  );
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.manifest.Manifest} returns this
 */
proto.manifest.Manifest.prototype.clearRoutesList = function () {
  return this.setRoutesList([]);
};

if (jspb.Message.GENERATE_TO_OBJECT) {
  /**
   * Creates an object representation of this proto.
   * Field names that are reserved in JavaScript and will be renamed to pb_name.
   * Optional fields that are not set will be set to undefined.
   * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
   * For the list of reserved names please see:
   *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
   * @param {boolean=} opt_includeInstance Deprecated. whether to include the
   *     JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @return {!Object}
   */
  proto.manifest.BuildInfo.prototype.toObject = function (opt_includeInstance) {
    return proto.manifest.BuildInfo.toObject(opt_includeInstance, this);
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.manifest.BuildInfo} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.manifest.BuildInfo.toObject = function (includeInstance, msg) {
    var f,
      obj = {
        maintainer: jspb.Message.getFieldWithDefault(msg, 1, ""),
        image: jspb.Message.getFieldWithDefault(msg, 2, ""),
        tag: jspb.Message.getFieldWithDefault(msg, 3, ""),
        path: jspb.Message.getFieldWithDefault(msg, 4, ""),
      };

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.manifest.BuildInfo}
 */
proto.manifest.BuildInfo.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.manifest.BuildInfo();
  return proto.manifest.BuildInfo.deserializeBinaryFromReader(msg, reader);
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.manifest.BuildInfo} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.manifest.BuildInfo}
 */
proto.manifest.BuildInfo.deserializeBinaryFromReader = function (msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
      case 1:
        var value = /** @type {string} */ (reader.readString());
        msg.setMaintainer(value);
        break;
      case 2:
        var value = /** @type {string} */ (reader.readString());
        msg.setImage(value);
        break;
      case 3:
        var value = /** @type {string} */ (reader.readString());
        msg.setTag(value);
        break;
      case 4:
        var value = /** @type {string} */ (reader.readString());
        msg.setPath(value);
        break;
      default:
        reader.skipField();
        break;
    }
  }
  return msg;
};

/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.manifest.BuildInfo.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.manifest.BuildInfo.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.manifest.BuildInfo} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.manifest.BuildInfo.serializeBinaryToWriter = function (message, writer) {
  var f = undefined;
  f = message.getMaintainer();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getImage();
  if (f.length > 0) {
    writer.writeString(2, f);
  }
  f = message.getTag();
  if (f.length > 0) {
    writer.writeString(3, f);
  }
  f = message.getPath();
  if (f.length > 0) {
    writer.writeString(4, f);
  }
};

/**
 * optional string maintainer = 1;
 * @return {string}
 */
proto.manifest.BuildInfo.prototype.getMaintainer = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};

/**
 * @param {string} value
 * @return {!proto.manifest.BuildInfo} returns this
 */
proto.manifest.BuildInfo.prototype.setMaintainer = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string image = 2;
 * @return {string}
 */
proto.manifest.BuildInfo.prototype.getImage = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};

/**
 * @param {string} value
 * @return {!proto.manifest.BuildInfo} returns this
 */
proto.manifest.BuildInfo.prototype.setImage = function (value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional string tag = 3;
 * @return {string}
 */
proto.manifest.BuildInfo.prototype.getTag = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};

/**
 * @param {string} value
 * @return {!proto.manifest.BuildInfo} returns this
 */
proto.manifest.BuildInfo.prototype.setTag = function (value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};

/**
 * optional string path = 4;
 * @return {string}
 */
proto.manifest.BuildInfo.prototype.getPath = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};

/**
 * @param {string} value
 * @return {!proto.manifest.BuildInfo} returns this
 */
proto.manifest.BuildInfo.prototype.setPath = function (value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};

if (jspb.Message.GENERATE_TO_OBJECT) {
  /**
   * Creates an object representation of this proto.
   * Field names that are reserved in JavaScript and will be renamed to pb_name.
   * Optional fields that are not set will be set to undefined.
   * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
   * For the list of reserved names please see:
   *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
   * @param {boolean=} opt_includeInstance Deprecated. whether to include the
   *     JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @return {!Object}
   */
  proto.manifest.Version.prototype.toObject = function (opt_includeInstance) {
    return proto.manifest.Version.toObject(opt_includeInstance, this);
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.manifest.Version} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.manifest.Version.toObject = function (includeInstance, msg) {
    var f,
      obj = {
        major: jspb.Message.getFieldWithDefault(msg, 1, 0),
        minor: jspb.Message.getFieldWithDefault(msg, 2, 0),
        patch: jspb.Message.getFieldWithDefault(msg, 3, 0),
      };

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.manifest.Version}
 */
proto.manifest.Version.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.manifest.Version();
  return proto.manifest.Version.deserializeBinaryFromReader(msg, reader);
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.manifest.Version} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.manifest.Version}
 */
proto.manifest.Version.deserializeBinaryFromReader = function (msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
      case 1:
        var value = /** @type {number} */ (reader.readInt32());
        msg.setMajor(value);
        break;
      case 2:
        var value = /** @type {number} */ (reader.readInt32());
        msg.setMinor(value);
        break;
      case 3:
        var value = /** @type {number} */ (reader.readInt32());
        msg.setPatch(value);
        break;
      default:
        reader.skipField();
        break;
    }
  }
  return msg;
};

/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.manifest.Version.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.manifest.Version.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.manifest.Version} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.manifest.Version.serializeBinaryToWriter = function (message, writer) {
  var f = undefined;
  f = message.getMajor();
  if (f !== 0) {
    writer.writeInt32(1, f);
  }
  f = message.getMinor();
  if (f !== 0) {
    writer.writeInt32(2, f);
  }
  f = message.getPatch();
  if (f !== 0) {
    writer.writeInt32(3, f);
  }
};

/**
 * optional int32 major = 1;
 * @return {number}
 */
proto.manifest.Version.prototype.getMajor = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};

/**
 * @param {number} value
 * @return {!proto.manifest.Version} returns this
 */
proto.manifest.Version.prototype.setMajor = function (value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};

/**
 * optional int32 minor = 2;
 * @return {number}
 */
proto.manifest.Version.prototype.getMinor = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};

/**
 * @param {number} value
 * @return {!proto.manifest.Version} returns this
 */
proto.manifest.Version.prototype.setMinor = function (value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};

/**
 * optional int32 patch = 3;
 * @return {number}
 */
proto.manifest.Version.prototype.getPatch = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};

/**
 * @param {number} value
 * @return {!proto.manifest.Version} returns this
 */
proto.manifest.Version.prototype.setPatch = function (value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};

if (jspb.Message.GENERATE_TO_OBJECT) {
  /**
   * Creates an object representation of this proto.
   * Field names that are reserved in JavaScript and will be renamed to pb_name.
   * Optional fields that are not set will be set to undefined.
   * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
   * For the list of reserved names please see:
   *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
   * @param {boolean=} opt_includeInstance Deprecated. whether to include the
   *     JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @return {!Object}
   */
  proto.manifest.Route.prototype.toObject = function (opt_includeInstance) {
    return proto.manifest.Route.toObject(opt_includeInstance, this);
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.manifest.Route} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.manifest.Route.toObject = function (includeInstance, msg) {
    var f,
      obj = {
        httpMethod: jspb.Message.getFieldWithDefault(msg, 1, 0),
        handler: jspb.Message.getFieldWithDefault(msg, 2, ""),
        pathSpec: jspb.Message.getFieldWithDefault(msg, 3, ""),
      };

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.manifest.Route}
 */
proto.manifest.Route.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.manifest.Route();
  return proto.manifest.Route.deserializeBinaryFromReader(msg, reader);
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.manifest.Route} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.manifest.Route}
 */
proto.manifest.Route.deserializeBinaryFromReader = function (msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
      case 1:
        var value = /** @type {!proto.manifest.HttpMethod} */ (reader.readEnum());
        msg.setHttpMethod(value);
        break;
      case 2:
        var value = /** @type {string} */ (reader.readString());
        msg.setHandler(value);
        break;
      case 3:
        var value = /** @type {string} */ (reader.readString());
        msg.setPathSpec(value);
        break;
      default:
        reader.skipField();
        break;
    }
  }
  return msg;
};

/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.manifest.Route.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.manifest.Route.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.manifest.Route} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.manifest.Route.serializeBinaryToWriter = function (message, writer) {
  var f = undefined;
  f = message.getHttpMethod();
  if (f !== 0.0) {
    writer.writeEnum(1, f);
  }
  f = message.getHandler();
  if (f.length > 0) {
    writer.writeString(2, f);
  }
  f = message.getPathSpec();
  if (f.length > 0) {
    writer.writeString(3, f);
  }
};

/**
 * optional HttpMethod http_method = 1;
 * @return {!proto.manifest.HttpMethod}
 */
proto.manifest.Route.prototype.getHttpMethod = function () {
  return /** @type {!proto.manifest.HttpMethod} */ (jspb.Message.getFieldWithDefault(
    this,
    1,
    0
  ));
};

/**
 * @param {!proto.manifest.HttpMethod} value
 * @return {!proto.manifest.Route} returns this
 */
proto.manifest.Route.prototype.setHttpMethod = function (value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};

/**
 * optional string handler = 2;
 * @return {string}
 */
proto.manifest.Route.prototype.getHandler = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};

/**
 * @param {string} value
 * @return {!proto.manifest.Route} returns this
 */
proto.manifest.Route.prototype.setHandler = function (value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional string path_spec = 3;
 * @return {string}
 */
proto.manifest.Route.prototype.getPathSpec = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};

/**
 * @param {string} value
 * @return {!proto.manifest.Route} returns this
 */
proto.manifest.Route.prototype.setPathSpec = function (value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};

/**
 * @enum {number}
 */
proto.manifest.HttpMethod = {
  GET: 0,
  PUT: 1,
  POST: 2,
  DELETE: 3,
};

goog.object.extend(exports, proto.manifest);