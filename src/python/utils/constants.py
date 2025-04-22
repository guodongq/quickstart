from types import SimpleNamespace


class Constants(SimpleNamespace):
    def __setattr__(self, key, value):
        if hasattr(self, key):
            raise AttributeError(f"Cannot reassign constant '{key}'")
        super().__setattr__(key, value)
