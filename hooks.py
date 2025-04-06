from datetime import datetime


def on_config(config, **kwargs):
    # Update copyright automatically.
    config.copyright = f"💚 {datetime.now().year} Created by jrosco."
