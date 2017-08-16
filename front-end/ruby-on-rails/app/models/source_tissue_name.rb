class SourceTissueName < ApplicationRecord
	belongs_to :tissue
	belongs_to :source
end
