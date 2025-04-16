from types import SimpleNamespace


class Constants(SimpleNamespace):
    def __setattr__(self, key, value):
        if hasattr(self, key):
            raise AttributeError(f"Cannot reassign constant '{key}'")
        super().__setattr__(key, value)


constants = Constants()
constants.BASE_URL = "https://api.openweathermap.org/data/2.5/weather"
constants.API_KEY = "bd5e378503939ddaee76f12ad7a97608"
constants.USER_AGENT = "weather-app/1.0"
