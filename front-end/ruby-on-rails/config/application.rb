require_relative 'boot'

require 'csv'
require 'rails/all'

# Require the gems listed in Gemfile, including any gems
# you've limited to :test, :development, or :production.
Bundler.require(*Rails.groups)

module Pharmacodb
  class Application < Rails::Application
    # Settings in config/environments/* take precedence over those specified here.
    # Application configuration should go into files in config/initializers
    # -- all .rb files in that directory are automatically loaded.

    config.exceptions_app = self.routes
    config.assets.paths << Rails.root.join("app", "assets", "files")
    config.assets.paths << Rails.root.join("app", "assets", "fonts")
    config.my_search_split = /(?i)\s+(and)?\s*|\s+with?\s*|\s+plus?\s*|\s*[&\+,]+\s*|\s+/

  end
end
