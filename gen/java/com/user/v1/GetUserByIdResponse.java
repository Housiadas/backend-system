// Generated by the protocol buffer compiler.  DO NOT EDIT!
// NO CHECKED-IN PROTOBUF GENCODE
// source: user/v1/user.proto
// Protobuf Java Version: 4.30.0

package com.user.v1;

/**
 * Protobuf type {@code user.v1.GetUserByIdResponse}
 */
public final class GetUserByIdResponse extends
    com.google.protobuf.GeneratedMessage implements
    // @@protoc_insertion_point(message_implements:user.v1.GetUserByIdResponse)
    GetUserByIdResponseOrBuilder {
private static final long serialVersionUID = 0L;
  static {
    com.google.protobuf.RuntimeVersion.validateProtobufGencodeVersion(
      com.google.protobuf.RuntimeVersion.RuntimeDomain.PUBLIC,
      /* major= */ 4,
      /* minor= */ 30,
      /* patch= */ 0,
      /* suffix= */ "",
      GetUserByIdResponse.class.getName());
  }
  // Use GetUserByIdResponse.newBuilder() to construct.
  private GetUserByIdResponse(com.google.protobuf.GeneratedMessage.Builder<?> builder) {
    super(builder);
  }
  private GetUserByIdResponse() {
  }

  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return com.user.v1.UserProto.internal_static_user_v1_GetUserByIdResponse_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessage.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return com.user.v1.UserProto.internal_static_user_v1_GetUserByIdResponse_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            com.user.v1.GetUserByIdResponse.class, com.user.v1.GetUserByIdResponse.Builder.class);
  }

  private int bitField0_;
  public static final int USER_FIELD_NUMBER = 1;
  private com.user.v1.User user_;
  /**
   * <code>.user.v1.User user = 1 [json_name = "user"];</code>
   * @return Whether the user field is set.
   */
  @java.lang.Override
  public boolean hasUser() {
    return ((bitField0_ & 0x00000001) != 0);
  }
  /**
   * <code>.user.v1.User user = 1 [json_name = "user"];</code>
   * @return The user.
   */
  @java.lang.Override
  public com.user.v1.User getUser() {
    return user_ == null ? com.user.v1.User.getDefaultInstance() : user_;
  }
  /**
   * <code>.user.v1.User user = 1 [json_name = "user"];</code>
   */
  @java.lang.Override
  public com.user.v1.UserOrBuilder getUserOrBuilder() {
    return user_ == null ? com.user.v1.User.getDefaultInstance() : user_;
  }

  private byte memoizedIsInitialized = -1;
  @java.lang.Override
  public final boolean isInitialized() {
    byte isInitialized = memoizedIsInitialized;
    if (isInitialized == 1) return true;
    if (isInitialized == 0) return false;

    memoizedIsInitialized = 1;
    return true;
  }

  @java.lang.Override
  public void writeTo(com.google.protobuf.CodedOutputStream output)
                      throws java.io.IOException {
    if (((bitField0_ & 0x00000001) != 0)) {
      output.writeMessage(1, getUser());
    }
    getUnknownFields().writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (((bitField0_ & 0x00000001) != 0)) {
      size += com.google.protobuf.CodedOutputStream
        .computeMessageSize(1, getUser());
    }
    size += getUnknownFields().getSerializedSize();
    memoizedSize = size;
    return size;
  }

  @java.lang.Override
  public boolean equals(final java.lang.Object obj) {
    if (obj == this) {
     return true;
    }
    if (!(obj instanceof com.user.v1.GetUserByIdResponse)) {
      return super.equals(obj);
    }
    com.user.v1.GetUserByIdResponse other = (com.user.v1.GetUserByIdResponse) obj;

    if (hasUser() != other.hasUser()) return false;
    if (hasUser()) {
      if (!getUser()
          .equals(other.getUser())) return false;
    }
    if (!getUnknownFields().equals(other.getUnknownFields())) return false;
    return true;
  }

  @java.lang.Override
  public int hashCode() {
    if (memoizedHashCode != 0) {
      return memoizedHashCode;
    }
    int hash = 41;
    hash = (19 * hash) + getDescriptor().hashCode();
    if (hasUser()) {
      hash = (37 * hash) + USER_FIELD_NUMBER;
      hash = (53 * hash) + getUser().hashCode();
    }
    hash = (29 * hash) + getUnknownFields().hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static com.user.v1.GetUserByIdResponse parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.user.v1.GetUserByIdResponse parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.user.v1.GetUserByIdResponse parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.user.v1.GetUserByIdResponse parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.user.v1.GetUserByIdResponse parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static com.user.v1.GetUserByIdResponse parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static com.user.v1.GetUserByIdResponse parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input);
  }
  public static com.user.v1.GetUserByIdResponse parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  public static com.user.v1.GetUserByIdResponse parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseDelimitedWithIOException(PARSER, input);
  }

  public static com.user.v1.GetUserByIdResponse parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static com.user.v1.GetUserByIdResponse parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input);
  }
  public static com.user.v1.GetUserByIdResponse parseFrom(
      com.google.protobuf.CodedInputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessage
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  @java.lang.Override
  public Builder newBuilderForType() { return newBuilder(); }
  public static Builder newBuilder() {
    return DEFAULT_INSTANCE.toBuilder();
  }
  public static Builder newBuilder(com.user.v1.GetUserByIdResponse prototype) {
    return DEFAULT_INSTANCE.toBuilder().mergeFrom(prototype);
  }
  @java.lang.Override
  public Builder toBuilder() {
    return this == DEFAULT_INSTANCE
        ? new Builder() : new Builder().mergeFrom(this);
  }

  @java.lang.Override
  protected Builder newBuilderForType(
      com.google.protobuf.GeneratedMessage.BuilderParent parent) {
    Builder builder = new Builder(parent);
    return builder;
  }
  /**
   * Protobuf type {@code user.v1.GetUserByIdResponse}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessage.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:user.v1.GetUserByIdResponse)
      com.user.v1.GetUserByIdResponseOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.user.v1.UserProto.internal_static_user_v1_GetUserByIdResponse_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessage.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.user.v1.UserProto.internal_static_user_v1_GetUserByIdResponse_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.user.v1.GetUserByIdResponse.class, com.user.v1.GetUserByIdResponse.Builder.class);
    }

    // Construct using com.user.v1.GetUserByIdResponse.newBuilder()
    private Builder() {
      maybeForceBuilderInitialization();
    }

    private Builder(
        com.google.protobuf.GeneratedMessage.BuilderParent parent) {
      super(parent);
      maybeForceBuilderInitialization();
    }
    private void maybeForceBuilderInitialization() {
      if (com.google.protobuf.GeneratedMessage
              .alwaysUseFieldBuilders) {
        internalGetUserFieldBuilder();
      }
    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      bitField0_ = 0;
      user_ = null;
      if (userBuilder_ != null) {
        userBuilder_.dispose();
        userBuilder_ = null;
      }
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return com.user.v1.UserProto.internal_static_user_v1_GetUserByIdResponse_descriptor;
    }

    @java.lang.Override
    public com.user.v1.GetUserByIdResponse getDefaultInstanceForType() {
      return com.user.v1.GetUserByIdResponse.getDefaultInstance();
    }

    @java.lang.Override
    public com.user.v1.GetUserByIdResponse build() {
      com.user.v1.GetUserByIdResponse result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public com.user.v1.GetUserByIdResponse buildPartial() {
      com.user.v1.GetUserByIdResponse result = new com.user.v1.GetUserByIdResponse(this);
      if (bitField0_ != 0) { buildPartial0(result); }
      onBuilt();
      return result;
    }

    private void buildPartial0(com.user.v1.GetUserByIdResponse result) {
      int from_bitField0_ = bitField0_;
      int to_bitField0_ = 0;
      if (((from_bitField0_ & 0x00000001) != 0)) {
        result.user_ = userBuilder_ == null
            ? user_
            : userBuilder_.build();
        to_bitField0_ |= 0x00000001;
      }
      result.bitField0_ |= to_bitField0_;
    }

    @java.lang.Override
    public Builder mergeFrom(com.google.protobuf.Message other) {
      if (other instanceof com.user.v1.GetUserByIdResponse) {
        return mergeFrom((com.user.v1.GetUserByIdResponse)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(com.user.v1.GetUserByIdResponse other) {
      if (other == com.user.v1.GetUserByIdResponse.getDefaultInstance()) return this;
      if (other.hasUser()) {
        mergeUser(other.getUser());
      }
      this.mergeUnknownFields(other.getUnknownFields());
      onChanged();
      return this;
    }

    @java.lang.Override
    public final boolean isInitialized() {
      return true;
    }

    @java.lang.Override
    public Builder mergeFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      if (extensionRegistry == null) {
        throw new java.lang.NullPointerException();
      }
      try {
        boolean done = false;
        while (!done) {
          int tag = input.readTag();
          switch (tag) {
            case 0:
              done = true;
              break;
            case 10: {
              input.readMessage(
                  internalGetUserFieldBuilder().getBuilder(),
                  extensionRegistry);
              bitField0_ |= 0x00000001;
              break;
            } // case 10
            default: {
              if (!super.parseUnknownField(input, extensionRegistry, tag)) {
                done = true; // was an endgroup tag
              }
              break;
            } // default:
          } // switch (tag)
        } // while (!done)
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.unwrapIOException();
      } finally {
        onChanged();
      } // finally
      return this;
    }
    private int bitField0_;

    private com.user.v1.User user_;
    private com.google.protobuf.SingleFieldBuilder<
        com.user.v1.User, com.user.v1.User.Builder, com.user.v1.UserOrBuilder> userBuilder_;
    /**
     * <code>.user.v1.User user = 1 [json_name = "user"];</code>
     * @return Whether the user field is set.
     */
    public boolean hasUser() {
      return ((bitField0_ & 0x00000001) != 0);
    }
    /**
     * <code>.user.v1.User user = 1 [json_name = "user"];</code>
     * @return The user.
     */
    public com.user.v1.User getUser() {
      if (userBuilder_ == null) {
        return user_ == null ? com.user.v1.User.getDefaultInstance() : user_;
      } else {
        return userBuilder_.getMessage();
      }
    }
    /**
     * <code>.user.v1.User user = 1 [json_name = "user"];</code>
     */
    public Builder setUser(com.user.v1.User value) {
      if (userBuilder_ == null) {
        if (value == null) {
          throw new NullPointerException();
        }
        user_ = value;
      } else {
        userBuilder_.setMessage(value);
      }
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>.user.v1.User user = 1 [json_name = "user"];</code>
     */
    public Builder setUser(
        com.user.v1.User.Builder builderForValue) {
      if (userBuilder_ == null) {
        user_ = builderForValue.build();
      } else {
        userBuilder_.setMessage(builderForValue.build());
      }
      bitField0_ |= 0x00000001;
      onChanged();
      return this;
    }
    /**
     * <code>.user.v1.User user = 1 [json_name = "user"];</code>
     */
    public Builder mergeUser(com.user.v1.User value) {
      if (userBuilder_ == null) {
        if (((bitField0_ & 0x00000001) != 0) &&
          user_ != null &&
          user_ != com.user.v1.User.getDefaultInstance()) {
          getUserBuilder().mergeFrom(value);
        } else {
          user_ = value;
        }
      } else {
        userBuilder_.mergeFrom(value);
      }
      if (user_ != null) {
        bitField0_ |= 0x00000001;
        onChanged();
      }
      return this;
    }
    /**
     * <code>.user.v1.User user = 1 [json_name = "user"];</code>
     */
    public Builder clearUser() {
      bitField0_ = (bitField0_ & ~0x00000001);
      user_ = null;
      if (userBuilder_ != null) {
        userBuilder_.dispose();
        userBuilder_ = null;
      }
      onChanged();
      return this;
    }
    /**
     * <code>.user.v1.User user = 1 [json_name = "user"];</code>
     */
    public com.user.v1.User.Builder getUserBuilder() {
      bitField0_ |= 0x00000001;
      onChanged();
      return internalGetUserFieldBuilder().getBuilder();
    }
    /**
     * <code>.user.v1.User user = 1 [json_name = "user"];</code>
     */
    public com.user.v1.UserOrBuilder getUserOrBuilder() {
      if (userBuilder_ != null) {
        return userBuilder_.getMessageOrBuilder();
      } else {
        return user_ == null ?
            com.user.v1.User.getDefaultInstance() : user_;
      }
    }
    /**
     * <code>.user.v1.User user = 1 [json_name = "user"];</code>
     */
    private com.google.protobuf.SingleFieldBuilder<
        com.user.v1.User, com.user.v1.User.Builder, com.user.v1.UserOrBuilder> 
        internalGetUserFieldBuilder() {
      if (userBuilder_ == null) {
        userBuilder_ = new com.google.protobuf.SingleFieldBuilder<
            com.user.v1.User, com.user.v1.User.Builder, com.user.v1.UserOrBuilder>(
                getUser(),
                getParentForChildren(),
                isClean());
        user_ = null;
      }
      return userBuilder_;
    }

    // @@protoc_insertion_point(builder_scope:user.v1.GetUserByIdResponse)
  }

  // @@protoc_insertion_point(class_scope:user.v1.GetUserByIdResponse)
  private static final com.user.v1.GetUserByIdResponse DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new com.user.v1.GetUserByIdResponse();
  }

  public static com.user.v1.GetUserByIdResponse getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<GetUserByIdResponse>
      PARSER = new com.google.protobuf.AbstractParser<GetUserByIdResponse>() {
    @java.lang.Override
    public GetUserByIdResponse parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      Builder builder = newBuilder();
      try {
        builder.mergeFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.setUnfinishedMessage(builder.buildPartial());
      } catch (com.google.protobuf.UninitializedMessageException e) {
        throw e.asInvalidProtocolBufferException().setUnfinishedMessage(builder.buildPartial());
      } catch (java.io.IOException e) {
        throw new com.google.protobuf.InvalidProtocolBufferException(e)
            .setUnfinishedMessage(builder.buildPartial());
      }
      return builder.buildPartial();
    }
  };

  public static com.google.protobuf.Parser<GetUserByIdResponse> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<GetUserByIdResponse> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public com.user.v1.GetUserByIdResponse getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

