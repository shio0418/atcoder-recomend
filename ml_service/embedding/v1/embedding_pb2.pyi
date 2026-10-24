from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class EmbedRequest(_message.Message):
    __slots__ = ("problem_id", "text")
    PROBLEM_ID_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    problem_id: str
    text: str
    def __init__(self, problem_id: _Optional[str] = ..., text: _Optional[str] = ...) -> None: ...

class EmbedResponse(_message.Message):
    __slots__ = ("problem_id", "values", "model_version")
    PROBLEM_ID_FIELD_NUMBER: _ClassVar[int]
    VALUES_FIELD_NUMBER: _ClassVar[int]
    MODEL_VERSION_FIELD_NUMBER: _ClassVar[int]
    problem_id: str
    values: _containers.RepeatedScalarFieldContainer[float]
    model_version: str
    def __init__(self, problem_id: _Optional[str] = ..., values: _Optional[_Iterable[float]] = ..., model_version: _Optional[str] = ...) -> None: ...
