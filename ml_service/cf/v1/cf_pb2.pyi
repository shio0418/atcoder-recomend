from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class RecommendRequest(_message.Message):
    __slots__ = ("user_handle", "k")
    USER_HANDLE_FIELD_NUMBER: _ClassVar[int]
    K_FIELD_NUMBER: _ClassVar[int]
    user_handle: str
    k: int
    def __init__(self, user_handle: _Optional[str] = ..., k: _Optional[int] = ...) -> None: ...

class RecommendResponse(_message.Message):
    __slots__ = ("problem_id", "model_version")
    PROBLEM_ID_FIELD_NUMBER: _ClassVar[int]
    MODEL_VERSION_FIELD_NUMBER: _ClassVar[int]
    problem_id: _containers.RepeatedScalarFieldContainer[str]
    model_version: str
    def __init__(self, problem_id: _Optional[_Iterable[str]] = ..., model_version: _Optional[str] = ...) -> None: ...
