from typing import Any
from pydantic import BaseModel, ConfigDict, field_validator
from pydantic_core import PydanticUndefined

class BaseSchema(BaseModel):
    model_config = ConfigDict (
        poulate_by_name = True,
        from_attributes = True,
        arbitrary_types_allowed=True
    )

    @field_validator("*", mode="before")
    @classmethod
    def use_default_value(cls, value: Any, info) -> Any:
        if (
            cls.model_fields[info.field_name].get_default is not PydanticUndefined
            and not cls.model_fields[info.field_name].is_required()
            and value is None
        ):
            return cls.model_fields[info.field_name].get_default()
        return value